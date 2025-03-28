package models

import play.api.libs.json._

case class CartDTO(productId: Int, quantity: Int)

object CartDTO {
  implicit val cartDTOFormat: OFormat[CartDTO] = Json.format[CartDTO]
}