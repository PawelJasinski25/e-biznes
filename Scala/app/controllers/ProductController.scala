package controllers

import models.Product
import models.ProductDTO
import play.api.mvc._
import play.api.libs.json._

import javax.inject._
import scala.collection.mutable

@Singleton
class ProductController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val products = mutable.ListBuffer(
    Product(1, "Laptop", 1, 3500.00),
    Product(2, "Telefon", 1, 2500.00),
    Product(3, "Tablet", 1, 1500.00),
    Product(4, "T-shirt", 2, 15.99),
    Product(5, "Jacket", 2, 85.00),
    Product(6, "Apple", 3, 2.50)

  )

  def getAll: Action[AnyContent] = Action {
    Ok(Json.toJson(products))
  }

  def getById(id: Int): Action[AnyContent] = Action {
    products.find(_.id == id) match {
      case Some(product) => Ok(Json.toJson(product))
      case None => NotFound(Json.obj("error" -> s"Product with ID $id not found"))
    }
  }

  def add: Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[ProductDTO].fold(
      _ => BadRequest(Json.obj("error" -> "Incorrect data")),
      productDTO => {
        val newId = if (products.isEmpty) 1 else products.map(_.id).max + 1
        val productWithId = Product(newId, productDTO.name, productDTO.categoryId, productDTO.price)
        products += productWithId
        Created(Json.toJson(productWithId))
      }
    )
  }

  def update(id: Int): Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[ProductDTO].fold(
      _ => BadRequest(Json.obj("error" -> "Incorrect data")),
      updatedProductDTO => {
        products.indexWhere(_.id == id) match {
          case -1 => NotFound(Json.obj("error" -> "Product not found"))
          case idx =>
            val updatedProduct = Product(id, updatedProductDTO.name, updatedProductDTO.categoryId, updatedProductDTO.price)
            products.update(idx, updatedProduct)
            Ok(Json.toJson(updatedProduct))
        }
      }
    )
  }

  def delete(id: Int): Action[AnyContent] = Action {
    products.indexWhere(_.id == id) match {
      case -1 => NotFound(Json.obj("error" -> "Product not found"))
      case idx =>
        products.remove(idx)
        NoContent
    }
  }
}