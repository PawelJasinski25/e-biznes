package models

import play.api.libs.json._

case class Product(id: Int, name: String, categoryId: Int, price: Double)

object Product {
  implicit val productFormat: OFormat[Product] = Json.format[Product]
}
