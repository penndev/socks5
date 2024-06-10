package com.penndev.socks5.ui.home


import android.os.Bundle
import android.text.InputType
import android.util.Log
import android.view.MenuItem
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import com.penndev.socks5.R
import com.penndev.socks5.databinding.ActivityHomeNodeBinding
import com.penndev.socks5.databinding.NodeData


class HomeNodeActivity : AppCompatActivity() {

    private lateinit var binding: ActivityHomeNodeBinding

    private lateinit var nodedata: NodeData

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityHomeNodeBinding.inflate(layoutInflater)
        nodedata = NodeData(this)
        setContentView(binding.root)
        initView()
    }

    private fun initView() {
        supportActionBar?.setDisplayHomeAsUpEnabled(true)
        // set old value.
        binding.inputType.setSelection(nodedata.type!!)
        binding.inputHost.setText(nodedata.host)
        nodedata.port?.let{if(it > 0)binding.inputPort.setText(it.toString())}
        binding.inputUser.setText(nodedata.user)
        binding.inputPass.setText(nodedata.pass)
        // set action
        binding.inputPassShow.setOnCheckedChangeListener { _, isChecked ->
            binding.inputPass.inputType = if (isChecked) InputType.TYPE_CLASS_TEXT
            else InputType.TYPE_CLASS_TEXT or InputType.TYPE_TEXT_VARIATION_PASSWORD
        }
        //Toast.makeText(this,"${nodedata.typeSelected} |tcp-> ${nodedata.typeTcpEnable} |udp-> ${nodedata.typeUdpEnable}", Toast.LENGTH_LONG).show()
        binding.inputSubmit.setOnClickListener {
            nodedata.type = binding.inputType.selectedItemPosition
            nodedata.host = binding.inputHost.text.toString()
            nodedata.port = binding.inputPort.text.toString().toIntOrNull()
            nodedata.user = binding.inputUser.text.toString()
            nodedata.pass = binding.inputPass.text.toString()
            if(nodedata.host == null || nodedata.host == ""){
                binding.inputHost.requestFocus()
                binding.inputHost.error = getString(R.string.home_node_input_not_null)
                return@setOnClickListener
            }
            if(nodedata.port == 0){
                binding.inputPort.requestFocus()
                binding.inputPort.error = getString(R.string.home_node_input_not_null)
                return@setOnClickListener
            }
            finish()
        }
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        return when (item.itemId) {
            android.R.id.home -> {
                finish()
                true
            }
            else ->  {
                super.onOptionsItemSelected(item)
            }
        }
    }

}