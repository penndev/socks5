package com.github.penndev;


import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.*;
import java.util.Arrays;

public class Service implements Runnable {
    private final Socket sock;
    private final OutputStream output;
    private final InputStream input;

    private final byte version = 0x05;
    private byte method = 0x00;
    private byte cmd;

    private String host;
    private int port;

    public Service(Socket sock) throws IOException {
        this.sock = sock;
        this.input = sock.getInputStream();
        this.output = sock.getOutputStream();
        // 确定当前支持的认证方法。
        if (Socks5.username != null) {
            method = 0x02;
        }
    }

    //https://datatracker.ietf.org/doc/html/rfc1928#section-3
    private void HandShake() throws Socks5Exception, IOException {
        var ver_n_method = input.readNBytes(2);
        if (ver_n_method[0] != version || ver_n_method[1] < 1) {
            output.write(new byte[]{0x05, (byte) 0xff});
            throw new Socks5Exception("socks5 version error[" + ver_n_method[0] + "]");
        }
        // - handshake
        var methods = input.readNBytes(ver_n_method[1]);

        boolean contains = false;
        for (byte i : methods) {
            if (i == method) {
                contains = true;
                break;
            }
        }
        if (contains) {
            output.write(new byte[]{0x05, method});
        } else {
            output.write(new byte[]{0x05, (byte) 0xff});
            throw new Socks5Exception("socks5 not match method[" + method + "]");
        }
        if (method == 0x02) { // https://datatracker.ietf.org/doc/html/rfc1929
            if (input.readNBytes(1)[0] != 1) {
                output.write(new byte[]{0x05, 0x01});
                throw new Socks5Exception("socks5 user pass VER error");
            }
            String username = new String(input.readNBytes(input.readNBytes(1)[0]));
            String password = new String(input.readNBytes(input.readNBytes(1)[0]));
            if (username.equals(Socks5.username) && password.equals(Socks5.password)) {
                output.write(new byte[]{0x05, 0x00});
            } else {
                output.write(new byte[]{0x05, (byte) 0xff});
                throw new Socks5Exception("socks5 user pass not right");
            }
        }
    }

    //https://datatracker.ietf.org/doc/html/rfc1928#section-4
    public void Requests() throws Socks5Exception, IOException {
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

    //https://datatracker.ietf.org/doc/html/rfc1928#section-6
    public void replies(byte[] ip, int port, byte rep) throws IOException {
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

    private void CmdConnect() throws IOException, Socks5Exception {
        Socket remote = new Socket();
        byte[] ip;
        //System.out.println(host + ":" + port);
        try {
            remote.setSoTimeout(30000);
            remote.connect(new InetSocketAddress(host, port));
            ip = remote.getInetAddress().getAddress();
        } catch (Exception e) {
            replies(new byte[]{0, 0, 0, 0}, port, (byte) 0x06);
            throw new Socks5Exception("Socks5 remote not connect");
        }
        if (remote.isConnected()) {
            replies(ip, port, (byte) 0x00);
            new Thread(() -> { // 远程传输给本地。
                try {
                    remote.getInputStream().transferTo(output);
                } catch (IOException e) {
                } finally {
                    if (!remote.isClosed()) try {
                        remote.close();
                    } catch (IOException e) {
                    }
                }
            }).start();

            try { // 本地传输给远程
                input.transferTo(remote.getOutputStream());
            } catch (IOException e) {
            }
        } else {
            replies(ip, port, (byte) 0x04);
        }
    }

    private void cmdUdp() throws IOException {
        Thread t;
        port = 0;
        host = sock.getLocalAddress().getHostAddress();
        byte[] lip = InetAddress.getByName(host).getAddress();
        try {
            UDPClient uc = new UDPClient();
            t = new Thread(uc, "UDP");
        } catch (Exception e) {
            // udp fail
            replies(lip, port, (byte) 0x03);
            throw e;
            //return;
        }
        t.start();
        replies(lip, port, (byte) 0x00);

        byte[] buffer = new byte[1024];
        int read;
        while ((read = input.read(buffer, 0, 1024)) >= 0) {
            byte[] slice = Arrays.copyOfRange(buffer, 0, read);
            System.out.println(Arrays.toString(slice));
        }
        t.interrupt();
    }

    @Override
    public void run() {
        try {
            HandShake();
            Requests();
            System.out.printf("CMD(%d) -> [%s:%d] \n", cmd, host, port);
            switch (cmd) {
                case 0x01 -> CmdConnect(); // CONNECT X'01'
                //case 0x02 -> // BIND X'02'
                case 0x03 -> cmdUdp(); // UDP ASSOCIATE X'03'
                default -> throw new Socks5Exception("Socks5 cmd Unexpected value:" + cmd);
            }
        } catch (Socks5Exception e) {
            // 正常审核错误
            System.out.println(e.getMessage());
        } catch (IOException e) {
            // 读写抛出错误
            System.out.println(e.getMessage());
        } catch (Exception e) {
            throw new RuntimeException(e);
        } finally {
            if (!sock.isClosed()) try {
                sock.close();
            } catch (IOException e) {
            }
        }
    }

    public class UDPClient implements Runnable {
        // udp 中继器
        private final DatagramSocket udpSock;
        //  客户端主机信息
        String srcAddr;
        int srcPort;
        // 目标主机信息
        String dstAddr;
        int dstPort;

        // ....
        public UDPClient() throws SocketException {
            //udpSock = new DatagramSocket(new InetSocketAddress(host, port));
            udpSock = new DatagramSocket(8181);
            port = udpSock.getLocalPort();
            System.out.printf("UDP listen [%s:%d] \n", host, port);
        }

        @Override
        public void run() {
            byte[] buffer = new byte[1500];
            try {
                while (true) {
                    DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                    udpSock.receive(packet);
                    byte[] uData = packet.getData();
                    // 不处理数据分片 FRAG
                    if (uData[2] != 0x00) {
                        break;
                    }

                    // 标注客户端连接
                    if (uData[0] == 0x00 && uData[1] == 0x00) {
                        byte[] payload;
                        switch (uData[3]) {
                            case 0x01: {
                                dstAddr = InetAddress.getByAddress(Arrays.copyOfRange(uData, 4, 8)).getHostAddress();
                                dstPort = ((uData[8] & 0xff) << 8) | (uData[9] & 0xff);
                                payload = Arrays.copyOfRange(uData, 10, packet.getLength() + 1);
                                break;
                            }
                            case 0x03: {
                                int domainToLen = uData[4] + 5;
                                dstAddr = new String(Arrays.copyOfRange(uData, 5, domainToLen));
                                dstPort = ((uData[domainToLen] & 0xff) << 8) | (uData[domainToLen + 1] & 0xff);
                                payload = Arrays.copyOfRange(uData, domainToLen + 1, packet.getLength() + 1);
                                break;
                            }
                            case 0x04: {
                                dstAddr = new String(Arrays.copyOfRange(uData, 4, 21));
                                dstPort = ((uData[21] & 0xff) << 8) | (uData[22] & 0xff);
                                payload = Arrays.copyOfRange(uData, 23, packet.getLength() + 1);
                                break;
                            }
                            default: {
                                throw new IOException("socks5 merge the host atyp:" + uData[3]);
                            }
                        }
                        udpSock.send(new DatagramPacket(
                                payload, payload.length,
                                InetAddress.getByName(dstAddr), dstPort
                        ));
                        srcAddr = packet.getAddress().getHostAddress();
                        srcPort = packet.getPort();
                        System.out.printf("[%s:%d] --> [%s:%d] \n", srcAddr, srcPort, dstAddr, dstPort);

                    } else {
                        String currentAddr = packet.getAddress().getHostAddress();
                        int currentPort = packet.getPort();

                        if (currentAddr == dstAddr && currentPort == dstPort) {
                            udpSock.send(new DatagramPacket(
                                    packet.getData(), packet.getLength(),
                                    InetAddress.getByName(srcAddr), srcPort
                            ));
                        } else if (currentAddr == srcAddr && currentPort == srcPort) {
                            udpSock.send(new DatagramPacket(
                                    packet.getData(), packet.getLength(),
                                    InetAddress.getByName(dstAddr), dstPort
                            ));
                        }
                        // 接收回复数据
                        System.out.println("UDP receive <- " + currentAddr + ":" + currentPort);
                    }

                }
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
    }
}
