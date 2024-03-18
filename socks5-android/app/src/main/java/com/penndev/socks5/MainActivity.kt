package com.penndev.socks5

import android.content.Intent
import android.net.VpnService
import android.os.Bundle
import android.view.View
import android.widget.EditText
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import com.penndev.socks5.service.Socks5Service


class MainActivity : AppCompatActivity() {

    val ALLOW_CREATE_SERVICE = 1

    private val intentSocks5Service: Intent
        get() = Intent(this, Socks5Service::class.java)

    private var inputHost: String
        get() = findViewById<EditText>(R.id.input_host).text.toString()
        set(value) = findViewById<EditText>(R.id.input_host).setText(value)

    private var inputPort: Int
        get() {
            return try {
                findViewById<EditText>(R.id.input_port).text.toString().toInt()
            } catch (e: NumberFormatException) { // 处理转换异常，例如返回默认端口号
                1080
            }
        }
        set(value) = findViewById<EditText>(R.id.input_port).setText(value.toString())

    private var inputUser: String
        get() = findViewById<EditText>(R.id.input_user).text.toString()
        set(value) = findViewById<EditText>(R.id.input_user).setText(value)

    private var inputPass: String
        get() = findViewById<EditText>(R.id.input_pass).text.toString()
        set(value) = findViewById<EditText>(R.id.input_pass).setText(value)

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
    }

    // 启动 Socks5Service的详细步骤和流程。
    private fun onStartSocks5Service() {
        VpnService.BIND_AUTO_CREATE
        val intentPrepare = VpnService.prepare(this)
        if (intentPrepare != null){
            startActivityForResult(intentPrepare,ALLOW_CREATE_SERVICE)
            return
        }
        if( inputHost == "" || inputPort < 1 ){
            Toast.makeText(this, "请输入服务器信息", Toast.LENGTH_LONG).show()
            return
        }

        val bundle = Bundle().apply {
            putString("host", inputHost)
            putInt("port", inputPort)
            putString("user", inputUser)
            putString("pass", inputPass)
        }
        startService(intentSocks5Service.putExtras(bundle))
    }

    fun handleAllowApp(view:View) {
        onStartAllowApp()
    }

    fun handleStart(view: View){
        onStartSocks5Service()
    }

    fun handleStop(view: View){
        //startService(intentService?.putExtra("close",true))
    }

    private fun onStartAllowApp() {
        val intent = Intent(this, AllowAppActivity::class.java)
        startActivity(intent)
    }

    // 请求用户授权开启vpn的结果
    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCode, resultCode, data)
        when(requestCode){
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