package models

import play.api.libs.json._

case class CategoryDTO(name: String)

object CategoryDTO {
  implicit val categoryDTOFormat: OFormat[CategoryDTO] = Json.format[CategoryDTO]
}