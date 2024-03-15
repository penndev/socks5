package com.github.penndev;

import java.io.IOException;
import java.net.ServerSocket;
import java.net.Socket;


public class Socks5 {
    /**
     * 监听端口
     */
    public static int Port = 1080;

    public static String username = null;

    public static String password = null;

    public static void main(String[] args) throws IOException {
        // 获取传递的参数
        if(args.length == 1 || args.length == 3){
            Socks5.Port = Integer.parseInt(args[0]);
            if(args.length == 3){
                Socks5.username = args[1];
                Socks5.password = args[2];
            }
        }else if(args.length  > 1){
            System.out.println("Error args:[port] [username] [password]");
            return;
        }

        System.out.println("Service start on:" + Socks5.Port);



        ServerSocket srvSocket = new ServerSocket(Socks5.Port);
        while (true) {
            Socket sock = srvSocket.accept();
            //System.out.println(sock.getRemoteSocketAddress());
            var t = new Thread(new Service(sock), "Service");
            t.start();
        }
    }
}
