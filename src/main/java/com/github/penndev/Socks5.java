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
        System.out.println("Socks5 listening on:" + Socks5.Port);
        ServerSocket srvSocket = new ServerSocket(Socks5.Port);
        while (true) {
            Socket sock = srvSocket.accept();
            //System.out.println(sock.getRemoteSocketAddress());
            var t = new Thread(new Service(sock), "Service");
            t.start();
        }
    }
}
