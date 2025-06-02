package com.example

import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.routing.*
import dev.kord.core.Kord
import dev.kord.common.entity.Snowflake
import dev.kord.core.event.message.MessageCreateEvent
import io.github.cdimascio.dotenv.dotenv
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import dev.kord.core.on
import com.slack.api.Slack
import com.slack.api.methods.SlackApiException
import com.slack.api.methods.request.chat.ChatPostMessageRequest
import com.slack.api.methods.response.chat.ChatPostMessageResponse
import io.ktor.client.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive
import io.ktor.client.plugins.HttpTimeout





suspend fun main() {
    val env = dotenv()
    val botToken = env["TOKEN"]
    val channelId = env["CHANNEL_ID"]
    val kordBot = Kord(botToken)

    val slackToken = env["SLACK_TOKEN"]
    val slackChannelId = env["SLACK_CHANNEL_ID"]
    val client = HttpClient() {
        install(HttpTimeout) {
            requestTimeoutMillis = 30000L
            connectTimeoutMillis = 15000L
            socketTimeoutMillis = 30000L
        }
    }
    val apiUrl = "http://127.0.0.1:8000/chat"


    val categories = listOf("games", "books", "electronics", "clothing", "sports")
    val productsByCategory = mapOf(
        "games" to listOf("Cyberpunk 2077", "The Witcher 3", "Red Dead Redemption 2"),
        "books" to listOf("Hobbit", "Harry Potter", "Dune"),
        "electronics" to listOf("Laptop", "Smartphone", "Headphones"),
        "clothing" to listOf("T-shirt", "Jeans", "Sneakers"),
        "sports" to listOf("Tennis Racket","Baseball Bat", "Soccer Ball")
    )

    fun sendMessageToSlack(message: String) {
        try {
            val slack = Slack.getInstance()
            val methods = slack.methods(slackToken)
            val response: ChatPostMessageResponse = methods.chatPostMessage(
                ChatPostMessageRequest.builder()
                    .channel(slackChannelId)
                    .text(message)
                    .build()
            )
            if (!response.isOk) {
                println("Slack Error: ${response.error}")
            } else {
                println("Message sent to Slack!")
            }
        } catch (e: SlackApiException) {
            println("Slack API Error: ${e.message}")
        } catch (e: Exception) {
            println("Error sending message to Slack: ${e.message}")
        }
    }

    embeddedServer(Netty, host = "127.0.0.1", port = 8080) {
        routing {
            post("/sendMessage") {
                val messageContent = call.receive<String>()
                kordBot.rest.channel.createMessage(Snowflake(channelId)) {
                    content = messageContent
                }
                call.respond(HttpStatusCode.OK, "Message sent successfully to channel with id $channelId.\"")
            }

            post("/sendMessageSlack") {
                val messageContent = call.receive<String>()
                sendMessageToSlack(messageContent)
                call.respond(HttpStatusCode.OK, "Message sent successfully to Slack channel with id $slackChannelId.")
            }
        }
    }.start(wait = false)

    kordBot.on<MessageCreateEvent> {
        val message = this.message
        if (!message.author?.isBot!!) {
            println("Received message from ${message.author?.username}: ${message.content}")
        }
        if (message.content == "!categories") {
            message.channel.createMessage("Categories: " + categories.joinToString(", "))
        }

        if (message.content.startsWith("!products")) {
            val parts = message.content.split(" ")
            if (parts.size > 1) {
                val category = parts[1].lowercase()
                val products = productsByCategory[category]
                message.channel.createMessage(
                    if (products != null)
                        "Products in the $category category: " + products.joinToString(", ")
                    else
                        "Unknown category"
                )
            } else {
                message.channel.createMessage("Please provide a category. Use: !products <category>")
            }
        }

        if (message.content.startsWith("!chat")) {
            val userMessage = message.content.removePrefix("!chat").trim()
            if (userMessage.isNotEmpty()) {
                try {
                    val response: HttpResponse = client.post(apiUrl) {
                        setBody("""{"user_message": "$userMessage"}""")
                        header(HttpHeaders.ContentType, ContentType.Application.Json)
                    }
                    val responseJson = Json.parseToJsonElement(response.bodyAsText()).jsonObject
                    val chatResponse = responseJson["response"]?.jsonPrimitive?.content ?: "Error retrieving response"

                    message.channel.createMessage(chatResponse)
                } catch (e: Exception) {
                    message.channel.createMessage("Error: Unable to connect to the chatbot service.")
                }
            } else {
                message.channel.createMessage("Please provide a message after !chat.")
            }
        }

    }

    kordBot.login()
}
