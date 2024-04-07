package com.github.penndev;


import org.jetbrains.annotations.NotNull;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.*;
import java.util.Arrays;
import java.util.concurrent.CompletableFuture;

@SuppressWarnings("FieldCanBeLocal")
public class Service implements Runnable {
    private final Socket sock;
    private final OutputStream output;
    private final InputStream input;

    private final byte version = 0x05;
    private byte method = 0x00;
    private byte cmd;

    private final byte methodNoAcceptable = (byte) 0xff;
    private final byte methodUsernamePassword = (byte) 0x02;

    private final byte replySucceeded = 0x00;
    private final byte replyNetworkUnreachable = 0x03;
    private final byte replyHostUnreachable = 0x04;

    private final int timeout = 5 * 1000;

    /**
     * proxy remote host and port
     */
    private String host;
    private int port;

    public Service(Socket sock) throws IOException {
        this.sock = sock;
        this.sock.setKeepAlive(true);
        this.input = sock.getInputStream();
        this.output = sock.getOutputStream();
        // 确定当前支持的认证方法。
        if (Socks5.username != null) {
            method = 0x02;
        }
    }

    @Override
    public void run() {
        try {
            handShake();
            requests();
            switch (cmd) {
                case 0x01 -> cmdConnect(); // CONNECT X'01'
                //case 0x02 -> // BIND X'02'
                case 0x03 -> cmdUdp(); // UDP ASSOCIATE X'03'
                default -> throw new Socks5Exception("Socks5 cmd Unexpected value:" + cmd);
            }
        } catch (Exception e) {
            try {
                sock.close();
            } catch (Exception ignore) {
            }
            System.out.println(e.getMessage());
        }
        //System.out.printf("Close from: [%s] \n", sock.getRemoteSocketAddress());
    }

    private boolean anyMatch(byte @NotNull [] methods, byte method) {
        boolean contains = false;
        for (byte i : methods) {
            if (i == method) {
                contains = true;
                break;
            }
        }
        return contains;
    }

    private void handShake() throws IOException {
        var verMethod = input.readNBytes(2);
        if (verMethod[0] != version || verMethod[1] < 1) {
            output.write(new byte[]{version, methodNoAcceptable});
            throw new Socks5Exception("socks5 version error[" + verMethod[0] + "]");
        }

        var methods = input.readNBytes(verMethod[1]);
        if (anyMatch(methods, method)) {
            output.write(new byte[]{version, method});
        } else {
            output.write(new byte[]{version, methodNoAcceptable});
            throw new Socks5Exception("socks5 not match method[" + method + "]");
        }
        //
        if (method == methodUsernamePassword) {
            if (input.readNBytes(1)[0] != 0x01) {
                output.write(new byte[]{0x01, 0x01});
                throw new Socks5Exception("socks5 user pass VER error");
            }
            String username = new String(input.readNBytes(input.readNBytes(1)[0]));
            String password = new String(input.readNBytes(input.readNBytes(1)[0]));
            if (username.equals(Socks5.username) && password.equals(Socks5.password)) {
                output.write(new byte[]{0x01, 0x00});
            } else {
                output.write(new byte[]{0x01, 0x01});
                throw new Socks5Exception("socks5 user pass not right");
            }
        }
    }

    private void requests() throws IOException {
        var d = input.readNBytes(4);
        if (d[0] != version || d.length != 4) {
            throw new Socks5Exception("socks5 version error 01");
        }
        cmd = d[1];
        host = switch (d[3]) {
            case 0x01 -> InetAddress.getByAddress(input.readNBytes(4)).getHostAddress();
            case 0x03 -> new String(input.readNBytes(input.readNBytes(1)[0]));
            case 0x04 -> new String(input.readNBytes(16));
            default -> throw new Socks5Exception("socks5 merge the host atyp:" + d[3]);
        };
        var portByte = input.readNBytes(2);
        port = (portByte[0] & 0xff) * 256 + (portByte[1] & 0xff);
    }

    public void replies(byte @NotNull [] ip, int port, byte rep) throws IOException {
        ByteArrayOutputStream replies = new ByteArrayOutputStream();
        replies.write(new byte[]{version, rep, 0x00});
        if (ip.length == 4) {
            replies.write(0x01);
            replies.write(ip);
        } else if (ip.length == 16) {
            replies.write(0x04);
            replies.write(ip);
        } else {
            replies.write(new byte[]{0x01, 0x00, 0x00, 0x00, 0x00});
        }
        replies.write(new byte[]{(byte) (port / 256), (byte) (port % 256)});
        output.write(replies.toByteArray());
    }

    public void tunnel( InputStream input, OutputStream out) throws IOException {
        var bufferSize = 8192;
        byte[] buffer = new byte[bufferSize];
        int read;
        while ((read = input.read(buffer, 0, bufferSize)) >= 0) {
            out.write(buffer, 0, read);
            out.flush();
        }
        throw new Socks5Exception("input close");
    }

    private void cmdConnect() throws IOException {
        System.out.printf("CMD(%d) -> [%s:%d] \n", cmd, host, port);
        Socket remote = new Socket();
        remote.setSoTimeout(timeout);
        remote.setKeepAlive(true);
        remote.connect(new InetSocketAddress(host, port));
        byte[] ip = remote.getInetAddress().getAddress();
        Runnable close = () -> {
            try {
                remote.close();
            } catch (IOException ignore) {
            }
            try {
                sock.close();
            } catch (IOException ignore) {
            }
        };
        if (remote.isConnected()) {
            replies(ip, port, replySucceeded);
            CompletableFuture.runAsync(() -> {
                try {
                    tunnel(remote.getInputStream(), output);
                } catch (IOException e) {
                    close.run();
                }
            });
            try {
                tunnel(input, remote.getOutputStream());
            } catch (IOException e) {
                close.run();
            }
        } else {
            replies(ip, port, replyHostUnreachable);
            close.run();
        }
    }

    private void cmdUdp() throws IOException {
        Thread t;
        UDPClient udpClient;
        try {
            port = 0;
            host = sock.getLocalAddress().getHostAddress();
            udpClient = new UDPClient();
            t = new Thread(udpClient);
            t.start();
        } catch (Exception e) {
            var hostAddress = InetAddress.getByName(host).getAddress();
            replies(hostAddress, port, replyNetworkUnreachable);
            throw e;
        }
        var hostAddress = InetAddress.getByName(host).getAddress();
        replies(hostAddress, port, replySucceeded);

        var buffer = new byte[16];
        //noinspection StatementWithEmptyBody
        while ((input.read(buffer, 0, 16)) >= 0) {
        }
        udpClient.close();
    }


    public class UDPClient implements Runnable {

        private final DatagramSocket udpSock; // udp 中继器

        private String srcAddress; //  客户端主机信息
        private int srcPort;

        private String dstAddress; // 目标主机信息
        private int dstPort;

        public UDPClient() throws SocketException {
            udpSock = new DatagramSocket(new InetSocketAddress(host, port));
            port = udpSock.getLocalPort();
            //System.out.printf("UDP listen [%s:%d] \n", host, port);
        }

        public void close() {
            try {
                udpSock.close();
            } catch (Exception ignore) {
            }
            try {
                sock.close();
            } catch (Exception ignore) {
            }
        }

        @Override
        public void run() {
            try {
                //noinspection InfiniteLoopStatement
                while (true) {
                    byte[] buffer = new byte[1024];
                    DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                    udpSock.receive(packet);
                    byte[] packetData = packet.getData();
                    int packetLen = packet.getLength();

                    if (packetData[0] == 0x00 && packetData[1] == 0x00) {
                        if (packetData[2] != 0x00) continue; // TODO FRAG
                        int fromLen = decodePacketHeader(packetData);
                        byte[] payload = Arrays.copyOfRange(packetData, fromLen, packetLen);
                        udpSock.send(new DatagramPacket(payload, payload.length,
                                InetAddress.getByName(dstAddress), dstPort));
                        srcAddress = packet.getAddress().getHostAddress();
                        srcPort = packet.getPort();
                        //System.out.printf("[%s:%d] --> [%s:%d] \n", srcAddress, srcPort, dstAddress, dstPort);
                    } else {
                        String currentAddress = packet.getAddress().getHostAddress();
                        int currentPort = packet.getPort();
                        if (currentAddress.equals(dstAddress) && currentPort == dstPort) {
                            byte[] packetHeader = encodePacketHeader(dstAddress, dstPort);
                            byte[] result = new byte[packetHeader.length + packetLen];
                            System.arraycopy(packetHeader, 0, result, 0, packetHeader.length);
                            System.arraycopy(packetData, 0, result, packetHeader.length, packetLen);

                            udpSock.send(new DatagramPacket(result, result.length,
                                    InetAddress.getByName(srcAddress), srcPort));

                            //System.out.println("UDP receive <- " + currentAddress + ":" + currentPort);
                        } else if (currentAddress.equals(srcAddress) && currentPort == srcPort) {
                            // 单次信息没传输完
                            udpSock.send(new DatagramPacket(packet.getData(), packet.getLength(),
                                    InetAddress.getByName(dstAddress), dstPort));
                        }
                    }
                }
            } catch (Exception ignore) {
            }
        }

        private int decodePacketHeader(byte[] packetData) throws IOException {
            switch (packetData[3]) {
                case 0x01 -> {
                    dstAddress = InetAddress.getByAddress(
                            Arrays.copyOfRange(packetData, 4, 8)).getHostAddress();
                    dstPort = ((packetData[8] & 0xff) << 8) | (packetData[9] & 0xff);
                    return 10;
                }
                case 0x03 -> {
                    int domainToLen = packetData[4] + 5;
                    dstAddress = new String(Arrays.copyOfRange(packetData, 5, domainToLen));
                    dstPort = ((packetData[domainToLen] & 0xff) << 8) | (packetData[domainToLen + 1] & 0xff);
                    return domainToLen + 2;
                }
                case 0x04 -> {
                    dstAddress = new String(Arrays.copyOfRange(packetData, 4, 21));
                    dstPort = ((packetData[21] & 0xff) << 8) | (packetData[22] & 0xff);
                    return 23;
                }
                default -> throw new IOException("socks5 merge the host atyp:" + packetData[3]);
            }
        }

        private byte[] encodePacketHeader(String address, int port) throws IOException {
            ByteArrayOutputStream replies = new ByteArrayOutputStream();
            replies.write(new byte[]{0x00, 0x00, 0x00});
            byte[] ip = InetAddress.getByName(address).getAddress();
            if (ip.length == 4) {
                replies.write(0x01);
                replies.write(ip);
            } else if (ip.length == 16) {
                replies.write(0x04);
                replies.write(ip);
            } else {
                replies.write(new byte[]{0x01, 0x00, 0x00, 0x00, 0x00});
            }
            replies.write(new byte[]{(byte) (port / 256), (byte) (port % 256)});
            return replies.toByteArray();
        }

    }
}
