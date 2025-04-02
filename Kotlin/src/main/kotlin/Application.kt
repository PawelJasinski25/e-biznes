package com.example

import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.routing.*
import dev.kord.core.Kord
import dev.kord.common.entity.Snowflake
import io.github.cdimascio.dotenv.dotenv
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*

suspend fun main() {
    val env = dotenv()
    val botToken = env["TOKEN"]
    val channelId = env["CHANNEL_ID"]
    val kordBot = Kord(botToken)

    embeddedServer(Netty, host = "127.0.0.1", port = 8080) {
        routing {
            post("/sendMessage") {
                val messageContent = call.receive<String>()
                kordBot.rest.channel.createMessage(Snowflake(channelId)) {
                    content = messageContent
                }
                call.respond(HttpStatusCode.OK, "Message sent successfully to channel with id $channelId.\"")
            }
        }
    }.start(wait = false)
    kordBot.login()
}
