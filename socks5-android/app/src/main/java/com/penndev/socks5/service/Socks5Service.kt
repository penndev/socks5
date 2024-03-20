package com.penndev.socks5.service

import android.annotation.SuppressLint
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Intent
import android.net.VpnService
import android.os.Build
import android.os.ParcelFileDescriptor
import android.util.Log
import android.widget.RemoteViews
import android.widget.Toast
import androidx.core.app.NotificationCompat
import com.penndev.socks5.MainActivity
import com.penndev.socks5.R
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.Job
import kotlinx.coroutines.launch
import java.net.InetAddress

class Socks5ServiceCloseException(message: String) : Exception(message)

class Socks5Service : VpnService() {

    // 回调UI状态
    interface OnStatus {
        fun start(){ status = true }
        fun close() { status = false }
    }

    companion object {
        var status:Boolean = false
        lateinit var onStatus: OnStatus
    }


    // tun 设备
    protected var tun: ParcelFileDescriptor? = null

    // 工作进程
    protected var job: Job? = null

    //远程服务器认证
    protected var serviceHost: String = ""
    protected var servicePort: Int = 0
    protected var serviceUser: String = ""
    protected var servicePass: String = ""
    protected var tunDNS: String = "8.8.8.8"
    protected var tunMtu: Int = 1400

    // 启动
    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        try {
            val err = setupCommand(intent)
            if (err != null) {
                Toast.makeText(this, err, Toast.LENGTH_SHORT).show()
            } else {
                setupNotifyForeground()
                onStatus.start()
            }
        }catch (e: Socks5ServiceCloseException) {
            onDestroy()
            Toast.makeText(this, e.message, Toast.LENGTH_SHORT).show()
        } catch (e: Exception) { //处理抛出异常问题
            onDestroy()
            Log.e("penndev", "启动异常", e)
        }
        return START_NOT_STICKY
    }

    // 启动初始化参数
    private fun setupCommand(intent: Intent?) : String? {
        if (intent == null) {
            return "传参错误"
        }
        if (intent.getBooleanExtra("close", false)) {
            throw Socks5ServiceCloseException("Closed")
        }
        if (status) {
            return "正在运行中"
        }

        serviceHost = intent.getStringExtra("host")!!

        try {
            InetAddress.getByName(serviceHost)
        }catch (e: Exception) {
            return "错误的host"
        }

        servicePort = intent.getIntExtra("port", 0)
        if (servicePort < 1 || servicePort > 65025){
            return "错误的port"
        }
        serviceUser =  intent.getStringExtra("user")!!
        servicePass = intent.getStringExtra("pass")!!
        return null
    }

    // 启动VPN
    fun setupVpnServe() {
        setupNotifyForeground() //启动通知

        tun = Builder().setMtu(tunMtu).addDnsServer(tunDNS)
            .addRoute("0.0.0.0", 0)
            .addAddress("192.168.1.1", 32)
            .addDisallowedApplication(packageName).establish()

        job = GlobalScope.launch {
            try {
                updateNotification("发送通知完成")
            } catch (e: Exception) { //处理抛出异常问题
                Log.e("penndev", "服务引起异常", e)
            }
        }
    }

    override fun onDestroy() {
        job?.cancel()
        tun?.close()
        onStatus.close()
        stopForeground(true)
        super.onDestroy()
    }

    override fun onRevoke() {
        //Toast.makeText(this, "TunVPN已停止", Toast.LENGTH_LONG).show()
        onDestroy()
    }

    //通知
    private val notifyID = 1
    private val notifyChannelID = "penndev.vpnService"
    private val notifyChannelName = "penndev.vpnService"
    private lateinit var notificationManager: NotificationManager
    private lateinit var notificationBuilder: NotificationCompat.Builder
    private lateinit var notificationView: RemoteViews
    @SuppressLint("UnspecifiedImmutableFlag")
    protected fun setupNotifyForeground() {
        notificationManager = getSystemService(NOTIFICATION_SERVICE) as NotificationManager

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(notifyChannelID, notifyChannelName, NotificationManager.IMPORTANCE_DEFAULT)
            notificationManager.createNotificationChannel(channel)
        }

        val intent = Intent(this, MainActivity::class.java)
        //val pendingIntent = PendingIntent.getActivity(this, 0, intent, PendingIntent.FLAG_UPDATE_CURRENT)
        val pendingIntent = PendingIntent.getActivity(this, 0, intent,
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE)


        notificationView = RemoteViews(packageName, R.layout.notification_layout)
        notificationView.setTextViewText(R.id.notificationTitle, "TunVPN")
        notificationView.setTextViewText(R.id.notificationContent, "正在连接服务器")

        notificationBuilder = NotificationCompat.Builder(this, notifyChannelID)
            .setSmallIcon(R.drawable.ic_launcher_foreground)
            .setCustomContentView(notificationView)
            .setContentIntent(pendingIntent)
            .setOngoing(true)
            .setOnlyAlertOnce(true)
            .setPriority(NotificationCompat.PRIORITY_DEFAULT)

        startForeground(notifyID, notificationBuilder.build())
    }

    protected fun updateNotification(contentText: String) {
        notificationView.setTextViewText(R.id.notificationContent, contentText)
        notificationManager.notify(notifyID, notificationBuilder.build())
    }
}

