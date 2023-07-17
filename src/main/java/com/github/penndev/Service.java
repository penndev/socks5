package com.github.penndev;


import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.*;

public class Service implements Runnable {
    private Socket sock;
    private OutputStream output;
    private InputStream input;

    private byte method = 0;

    private byte cmd;

    private String host;
    private int port;


    public Service(Socket sock) throws IOException {
        this.sock = sock;
        this.input = sock.getInputStream();
        this.output = sock.getOutputStream();
        // 确定当前支持的认证方法。
        if (Socks5.username != null) {
            method = 2;
        }

    }

    //https://datatracker.ietf.org/doc/html/rfc1928#section-3
    public boolean Connects() throws IOException {
        var ver_n_method = input.readNBytes(2);
        if (ver_n_method[0] != 0x05 || ver_n_method[1] < 1) {
            output.write(new byte[]{0x05, (byte) 0xff});
            return false;
        }

        var methods = input.readNBytes(ver_n_method[1]);
        for (byte i : methods) {
            if (i == method) {
                output.write(new byte[]{0x05, method});
                //https://datatracker.ietf.org/doc/html/rfc1929#section-2
                if (i == 0x02) {
                    if (input.readNBytes(1)[0] != 1) {
                        output.write(new byte[]{0x05, 0x01});
                        return false;
                    }
                    String username = new String(input.readNBytes(input.readNBytes(1)[0]));
                    String password = new String(input.readNBytes(input.readNBytes(1)[0]));
                    if(username.equals(Socks5.username) && password.equals(Socks5.password)){
                        output.write(new byte[]{0x05, 0x00});
                        return true;
                    }else{
                        output.write(new byte[]{0x05, (byte) 0xff});
                        return false;
                    }

                }
                return true;
            }
        }
        output.write(new byte[]{0x05, (byte) 0xff});
        return false;


    }

    //https://datatracker.ietf.org/doc/html/rfc1928#section-4
    public boolean Requests() throws IOException {
        var d = input.readNBytes(4);
        if (d[0] != 0x05) {
            return false;
        }
        cmd = d[1];
        host = switch (d[3]) {
            case 0x01 -> InetAddress.getByAddress(input.readNBytes(4)).getHostAddress();
            case 0x03 -> new String(input.readNBytes(input.readNBytes(1)[0]));
            case 0x04 -> new String(input.readNBytes(16));
            default -> throw new IllegalArgumentException("not merge the atyp ->" + d[3]);
        };
        var portByte = input.readNBytes(2);
        port = (portByte[0] & 0xff) * 256 + (portByte[1] & 0xff);
        return true;
    }

    //https://datatracker.ietf.org/doc/html/rfc1928#section-6
    public void Replies(byte[] ip, int port, byte rep) throws IOException {
        ByteArrayOutputStream replies = new ByteArrayOutputStream();
        replies.write(new byte[]{0x05, rep, 0x00, 0x04});
        replies.write(ip);
        replies.write(new byte[]{(byte) (port / 256), (byte) (port % 256)});
        output.write(replies.toByteArray());
    }

    private void cmdConnect() throws IOException {
        Socket remote = new Socket(host, port);
        if (remote.isConnected()) {
            Replies(remote.getInetAddress().getAddress(), port, (byte) 0x00);
            // 远程传输给本地。
            new Thread(() -> {
                try {
                    remote.getInputStream().transferTo(output);
                    if (!remote.isClosed()) remote.close();
                    if (!sock.isClosed()) sock.close();
                } catch (IOException e) {
                    return;
                }
            }).start();
            try {
                input.transferTo(remote.getOutputStream());
                if (!remote.isClosed()) remote.close();
                if (!sock.isClosed()) sock.close();
            } catch (IOException e) {
                return;
            }
        } else {
            Replies(remote.getInetAddress().getAddress(), port, (byte) 0x04);
        }
    }

    private void cmdUdp() throws SocketException {
        DatagramSocket remote = new DatagramSocket();
        new Thread(() -> {
            while (true) {
                byte[] receiveData = new byte[10240];
                DatagramPacket receivePacket = new DatagramPacket(receiveData, receiveData.length);
                try {
                    remote.receive(receivePacket);
                    output.write(receivePacket.getData(), 0, receivePacket.getLength());
                } catch (IOException e) {
                    return;
                }
            }
        }).start();

        try {
            // 本地写入远程。
            while (true) {
                byte[] sendData = new byte[0];
                sendData = input.readAllBytes();
                DatagramPacket sendPacket = new DatagramPacket(sendData, sendData.length, InetAddress.getByName(host), port);
                remote.send(sendPacket);
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

    }

    @Override
    public void run() {
        try {
            // 选择授权方式。
            if (!Connects()) {
                return;
            }
            if (!Requests()) {
                return;
            }
            System.out.println(host + ":" + port);
            switch (cmd) {
                case 0x01 -> cmdConnect(); // CONNECT X'01'
                //case 0x02 -> // BIND X'02'
                case 0x03 -> cmdUdp(); // UDP ASSOCIATE X'03'
                default -> throw new IllegalStateException("Unexpected value: " + cmd);
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
}
