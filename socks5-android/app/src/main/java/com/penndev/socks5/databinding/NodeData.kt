package com.penndev.socks5.databinding

import android.content.Context
import android.content.SharedPreferences
import androidx.core.content.edit
import com.penndev.socks5.R

class NodeData(context: Context) {

    private val ctx = context

    private val sharedPreferences: SharedPreferences = context.applicationContext.getSharedPreferences(
        context.applicationContext.packageName, Context.MODE_PRIVATE)

    var type: Int?
        get() = sharedPreferences.getInt("inputType", 0)
        set(value) {
            sharedPreferences.edit {
                if (value != null) {
                    putInt("inputType", value)
                }
            }
        }

    var typeSelected: String? = null
        get() {
            val typePositon = type!!
            val nodeInputTypeArray = ctx.resources.getStringArray(R.array.nodeInputType)
            return if (typePositon in nodeInputTypeArray.indices) {
                nodeInputTypeArray[typePositon]
            } else {
                null
            }
        }


    var host: String?
        get() = sharedPreferences.getString("inputHost", null)
        set(value) {
            sharedPreferences.edit {
                putString("inputHost", value)
            }
        }

    var port: Int?
        get() = sharedPreferences.getInt("inputPort", 0)
        set(value) {
            sharedPreferences.edit {
                if (value != null) {
                    putInt("inputPort", value)
                }
            }
        }

    var user: String?
        get() = sharedPreferences.getString("inputUser", null)
        set(value) {
            sharedPreferences.edit {
                putString("inputUser", value)
            }
        }

    var pass: String?
        get() = sharedPreferences.getString("inputPass", "")
        set(value) {
            sharedPreferences.edit {
                putString("inputPass", value)
            }
        }
}
