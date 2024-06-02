package com.penndev.socks5

import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.os.Build
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.navigation.findNavController
import androidx.navigation.ui.setupWithNavController
import com.penndev.socks5.databinding.ActivityMainBinding
import com.penndev.socks5.service.Socks5Service

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        //val navView = binding.navView
        //val navController = findNavController(R.id.nav_host_fragment_activity_main)
        //navView.setupWithNavController(navController)
        //
        //initView()
    }

    private fun initView() {
        setTitle(R.string.app_name)
        createNotificationChannel()
    }

    /**
     * 应用启动首先创建消息通道
     */
    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val name = Socks5Service.notifyChannelName
            val descriptionText = getString(R.string.notify_channel_description)
            val importance = NotificationManager.IMPORTANCE_DEFAULT
            val channel = NotificationChannel(Socks5Service.notifyChannelID, name, importance).apply {
                description = descriptionText
            }
            // Register the channel with the system.
            val notificationManager: NotificationManager =
                getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
            notificationManager.createNotificationChannel(channel)
        }
    }

}
