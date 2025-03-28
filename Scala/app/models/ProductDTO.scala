package models

import play.api.libs.json._

case class ProductDTO(name: String, categoryId: Int, price: Double)

object ProductDTO {
  implicit val productDTOFormat: OFormat[ProductDTO] = Json.format[ProductDTO]
}