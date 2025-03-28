package controllers

import models.Category
import models.CategoryDTO
import play.api.mvc._
import play.api.libs.json._

import javax.inject._
import scala.collection.mutable

@Singleton
class CategoryController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val categories = mutable.ListBuffer(
    Category(1, "Electronics"),
    Category(2, "Clothing"),
    Category(3, "Groceries")

  )

  def getAll: Action[AnyContent] = Action {
    Ok(Json.toJson(categories))
  }

  def getById(id: Int): Action[AnyContent] = Action {
    categories.find(_.id == id) match {
      case Some(categories) => Ok(Json.toJson(categories))
      case None => NotFound(Json.obj("error" -> s"Category with ID $id not found"))
    }
  }

  def add: Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[CategoryDTO].fold(
      _ => BadRequest(Json.obj("error" -> "Incorrect data")),
      categoryDTO => {
        val newId = if (categories.isEmpty) 1 else categories.map(_.id).max + 1
        val categoriesWithId = Category(newId, categoryDTO.name)
        categories += categoriesWithId
        Created(Json.toJson(categoriesWithId))
      }
    )
  }

  def update(id: Int): Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[CategoryDTO].fold(
      _ => BadRequest(Json.obj("error" -> "Incorrect data")),
      updatedCategoryDTO => {
        categories.indexWhere(_.id == id) match {
          case -1 => NotFound(Json.obj("error" -> "Category not found"))
          case idx =>
            val updatedCategory = Category(id, updatedCategoryDTO.name)
            categories.update(idx, updatedCategory)
            Ok(Json.toJson(updatedCategory))
        }
      }
    )
  }

  def delete(id: Int): Action[AnyContent] = Action {
    categories.indexWhere(_.id == id) match {
      case -1 => NotFound(Json.obj("error" -> "Category not found"))
      case idx =>
        categories.remove(idx)
        NoContent
    }
  }
}