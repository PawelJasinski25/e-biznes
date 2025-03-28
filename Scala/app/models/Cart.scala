package models

import play.api.libs.json._

case class Cart(id: Int, products: List[CartItem])

case class CartItem(productId: Int, quantity: Int)

object Cart {
  implicit val cartItemFormat: OFormat[CartItem] = Json.format[CartItem]
  implicit val cartFormat: OFormat[Cart] = Json.format[Cart]
}
