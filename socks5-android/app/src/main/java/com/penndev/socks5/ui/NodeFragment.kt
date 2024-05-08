package com.penndev.socks5.ui

import android.content.Context
import android.content.SharedPreferences
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.EditText
import androidx.core.content.edit
import androidx.fragment.app.Fragment
import com.google.android.material.bottomsheet.BottomSheetDialog
import com.penndev.socks5.R
import com.penndev.socks5.databinding.FragmentNodeBinding
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.InputStream
import java.io.OutputStream
import java.net.InetSocketAddress
import java.net.Socket


class NodeFragment : Fragment() {

    private lateinit var binding: FragmentNodeBinding


    private var sharedPreferences: SharedPreferences? = null
        get() {
            var ctx = requireContext()
            return ctx.getSharedPreferences(ctx.packageName, Context.MODE_PRIVATE)
        }

    var Host: String?
        get() = sharedPreferences?.getString("inputHost", "")
        set(value) {
            sharedPreferences?.edit(true) {
                putString("inputHost", value)
            }
        }
    var Port: String?
        get() = sharedPreferences?.getString("inputPort", "")
        set(value) {
            sharedPreferences?.edit(true) {
                putString("inputPort", value)
            }
        }

    var user: String? = ""
        get() = sharedPreferences?.getString("inputUser", "")

    var pass: String? = ""
        get() = sharedPreferences?.getString("inputPass", "")

    override fun onCreateView(inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {
        binding = FragmentNodeBinding.inflate(inflater, container, false)
        initView()
        return binding.root
    }

    private fun initView() {
        val bottomSheetDialog = BottomSheetDialog(requireContext())
        val bottomSheetView = layoutInflater.inflate(R.layout.dialog_node, null)
        bottomSheetDialog.setContentView(bottomSheetView)
        bottomSheetView.findViewById<Button>(R.id.input_submit).setOnClickListener {
            Host = bottomSheetView.findViewById<EditText>(R.id.input_host).text.toString()
            Port = bottomSheetView.findViewById<EditText>(R.id.input_port).text.toString()
            bottomSheetDialog.hide()
            setNodeSpeed()
        }
        binding.root.setOnClickListener{
            if (Host != null) bottomSheetView.findViewById<EditText>(R.id.input_host)?.setText(Host)
            if (Port != null) bottomSheetView.findViewById<EditText>(R.id.input_port)?.setText(Port)
            if (user != null) bottomSheetView.findViewById<EditText>(R.id.input_user)?.setText(user)
            if (pass != null) bottomSheetView.findViewById<EditText>(R.id.input_pass)?.setText(pass)
            bottomSheetDialog.show()
        }
        setNodeSpeed()
    }

    private fun setNodeSpeed() {
        if(Host?.length!! > 0 && Port?.length!! > 0){
            binding.currentNodeHost.text = "${Host}:${Port}"
            CoroutineScope(Dispatchers.IO).launch {
                var ms = hostSocks5RTT(Host!!, Port!!.toInt())
                withContext(Dispatchers.Main) {
                    binding.currentNodeSpeed.text = "${ms}ms"
                }
            }
        }
    }
    private fun hostSocks5RTT(host:String,port:Int): Int {
        try {
            return Socket().use { socket ->
                val timeout = 10000
                socket.connect(InetSocketAddress(host, port), timeout)
                socket.soTimeout = timeout
                val outputStream: OutputStream = socket.getOutputStream()
                val inputStream: InputStream = socket.getInputStream()

                val startTime = System.currentTimeMillis()
                outputStream.write(byteArrayOf(0x05, 0x01, 0x00)); outputStream.flush() // 无密码探测

                val packet = ByteArray(2)
                val bytesRead = inputStream.read(packet)
                val endTime = System.currentTimeMillis()

                if (bytesRead == 2 && packet[0].toInt() == 0x05) {
                    return (endTime - startTime).toInt()
                } else {
                    return 0
                }
            }
        } catch (e: Exception) {
            return 0
        }
        return 0
    }

}
