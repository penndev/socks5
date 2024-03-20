package com.penndev.socks5

import android.content.Intent
import android.net.VpnService
import android.os.Bundle
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import com.penndev.socks5.databinding.ActivityMainBinding
import com.penndev.socks5.service.Socks5Service

class MainActivity : AppCompatActivity() {
    //获取vpn启动权限标志
    private val ALLOW_CREATE_SERVICE = 1

    //view binding实例
    private lateinit var binding: ActivityMainBinding

    // 绑定UI和Socks5Service状态
    fun onStatus(): Socks5Service.OnStatus {
        return object : Socks5Service.OnStatus {
            override fun start() {
                super.start()
                binding.handleActionIcon.setImageResource(R.drawable.activity_main_start_close)
                binding.handleActionText.setText(R.string.activity_main_start_close_text)
            }

            override fun close() {
                super.close()
                binding.handleActionIcon.setImageResource(R.drawable.activity_main_start)
                binding.handleActionText.setText(R.string.activity_main_start_text)
            }
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        //设置各种初始化
        binding.handleActionIcon.setOnClickListener{
            if (Socks5Service.status)  onStopSocks5Service() else onStartSocks5Service()
        }
        Socks5Service.onStatus = onStatus()
    }

    private fun onStopSocks5Service() {
        startService(Intent(this, Socks5Service::class.java).putExtra("close",true))
    }

    // 启动 Socks5Service的详细步骤和流程。
    private fun onStartSocks5Service() {
        val intentPrepare = VpnService.prepare(this)
        if (intentPrepare != null) {
            startActivityForResult(intentPrepare, ALLOW_CREATE_SERVICE)
            return
        }

        val inputHost = binding.inputHost.text.toString()
        val inputPortStr = binding.inputPort.text.toString()
        val inputUser = binding.inputUser.text.toString()
        val inputPass = binding.inputPass.text.toString()
        if (inputHost == "" || inputPortStr == "") {
            Toast.makeText(this, getString(R.string.toast_main_param_error), Toast.LENGTH_LONG).show()
            return
        }
        val inputPort = inputPortStr.toInt()

        val bundle = Bundle().apply {
            putString("host", inputHost)
            putInt("port", inputPort)
            putString("user", inputUser)
            putString("pass", inputPass)
        }
        startService(Intent(this, Socks5Service::class.java).putExtras(bundle))
    }

    // 请求用户授权开启vpn的结果
    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCode, resultCode, data)
        when (requestCode) {
            ALLOW_CREATE_SERVICE -> {
                if (resultCode == RESULT_OK) {
                    onStartSocks5Service()
                } else {
                    Toast.makeText(this, "您拒绝了VPN请求[$resultCode]", Toast.LENGTH_LONG).show()
                }
            }
        }
    }
}
