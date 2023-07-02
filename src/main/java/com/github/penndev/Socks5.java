package com.github.penndev;

import java.io.*;
import java.net.ServerSocket;
import java.net.Socket;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

public class Socks5 {
    /**
     * 监听端口
     */
    public static final int port = 1080;

    /**
     *   认证方法
     *o  X'00' NO AUTHENTICATION REQUIRED
     *o  X'01' GSSAPI
     *o  X'02' USERNAME/PASSWORD
     *o  X'03' to X'7F' IANA ASSIGNED
     *o  X'80' to X'FE' RESERVED FOR PRIVATE METHODS
     *o  X'FF' NO ACCEPTABLE METHODS
     */
    public static final int method = 2;

    /**
     *  Socks Protocol version
     *  socks 版本首字符
     */
    public static final int VERSION = 5;

    /**
     * 预留字符默认值  socks5 默认为 0x00
     */
    public static final int RESERVED = 0;

    public static class Service implements Runnable {
        private Socket sock;
        private OutputStream output;
        private InputStream input;
        private String host;
        private Integer port;

        /**
         * o  CONNECT X'01'
         * o  BIND X'02'
         * o  UDP ASSOCIATE X'03'
         */
        private Integer cmd;

        /**
         * IP V4地址: 0x01
         * 域名地址: 0x03
         * IP V6地址: 0x04
         */
        private Integer type;

        /**
         * 请求地址类型 Addressing
         */
        public static final int TYPE_IPV4 = 1;
        public static final int TYPE_DOMAIN_NAME = 3;
        public static final int TYPE_IPV6 = 4;

        /**
         * 应答请求状态码
         */
        public static final byte repliesSuccess = 0; // 返回连接状态成功
        public static final int repliesNetworkUnreachable = 3; // 连接状态网络错误
        public static final byte repliesConnectionRefused = 5; // 远程服务器拒绝访问

        public Service(Socket sock) throws IOException {
            this.sock = sock;
            this.input = sock.getInputStream();
            this.output =  sock.getOutputStream();
        }

        private boolean usernamePassword() throws IOException {
            if( this.input.readNBytes(1)[0] != 1 ){
                return false;
            }
            int nameLen = this.input.readNBytes(1)[0] & 0xff;
            String username = new String(this.input.readNBytes(nameLen));
            int passLen = this.input.readNBytes(1)[0] & 0xff;
            String password = new String(this.input.readNBytes(passLen));
            System.out.println(username + ":" + password);
            return true;
        }

        public boolean authentication() throws IOException {
            byte[] head = new byte[256];
            int headLen = this.input.read(head);
            if (head[0] != VERSION ){
                return false;
            }
            for(int i = 1; i < headLen; i++){
                if (head[i] == method){
                    if (method == 0){
                        this.output.write(new byte[]{VERSION, 0});
                        return true;
                    }else if (method == 2){
                        this.output.write(new byte[]{VERSION, 2});
                        if(this.usernamePassword()){
                            this.output.write(new byte[]{VERSION, 0});
                            return true;
                        }else{
                            this.output.write(new byte[]{VERSION, 1});
                            return false;
                        }
                    }else {
                        throw new IOException("undefine auth method " + method);
                    }
                }
            }
            this.output.write(new byte[]{VERSION, (byte) 0xff});
            return false;
        }

        public boolean request() throws IOException {
            byte[] req = this.input.readNBytes(4);
            if (req[0] != 5) {
                System.out.println("协议版本错误");
                return false;
            }
            this.cmd = req[1] & 0xff;
            if (req[2] != 0) {
                System.out.println("保留字段错误");
                return false;
            }
            this.type = req[3] & 0xff;
            if (this.type == 0x01) {
                byte[] ipv4 = this.input.readNBytes(4);
                this.host = (ipv4[0] & 0xff) + "." + (ipv4[1] & 0xff) + "." + (ipv4[2] & 0xff) + "." + (ipv4[3] & 0xff);
            } else if (this.type == 0x03) {
                byte[] uriLen = this.input.readNBytes(1);
                this.host = new String(this.input.readNBytes(uriLen[0]));
            } else if (this.type == 0x04) {
                byte[] uriLen = this.input.readNBytes(16);
                this.host = new String(this.input.readNBytes(uriLen[0] & 0xff));
            } else {
                System.out.println("Addressing 错误的数据");
                return false;
            }
            byte[] portB16 = this.input.readNBytes(2);
            this.port = (portB16[0] & 0xff) * 256 + (portB16[1] & 0xff);
            return true;
        }

        private void replies(byte status) throws IOException {
            ByteArrayOutputStream replies = new ByteArrayOutputStream();
            replies.write(VERSION);
            replies.write(status);
            replies.write(new byte[]{RESERVED,TYPE_DOMAIN_NAME,});
            replies.write((byte)this.host.length());
            replies.write(this.host.getBytes());
            replies.write(new byte[]{(byte) (port / 256), (byte) (port % 256)});
            this.output.write(replies.toByteArray());
        }

        public void connect() throws IOException {
            Socket remote = new Socket(this.host,this.port);
            if (remote.isConnected()){
                this.replies( repliesSuccess);
                new Thread(()->{
                    try {
                        InputStream in = remote.getInputStream();
                        in.transferTo(this.output);
                        System.out.println(this.host + ":" + this.port + "-> remote closed!");
                    } catch (IOException e) {
                        throw new RuntimeException(e);
                    }
                }).start();
                OutputStream out = remote.getOutputStream();
                this.input.transferTo(out);
                System.out.println(this.host + ":" + this.port + "-> client closed!");
            }else{
                this.replies(repliesConnectionRefused);
            }
        }

        @Override
        public void run() {
            try {
                // 用户进行认证
                if (!this.authentication()) {
                    this.sock.close();
                    return;
                }
                // 开始数据传输
                if (!this.request()) {
                    this.sock.close();
                    return;
                }
                if (this.cmd == 0x1){
                    this.connect();
                }
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
    }

    public static void main(String[] args) throws IOException {
        System.out.println("Service start on :1080");
        ServerSocket srvSocket  = new ServerSocket(Socks5.port);
        while (true){
            Socket sock = srvSocket.accept();
            DateTimeFormatter formatStr = DateTimeFormatter.ofPattern("dd-MM-yyyy HH:mm:ss");
            String format = LocalDateTime.now().format(formatStr);
            System.out.println(format + " | new client " + sock.getRemoteSocketAddress());
            new Thread(new Service(sock),"listener").start();
        }
    }
}