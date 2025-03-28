package controllers

import models.Cart
import models.CartDTO
import models.CartItem
import play.api.mvc._
import play.api.libs.json._

import javax.inject._
import scala.collection.mutable

@Singleton
class CartController @Inject()(val controllerComponents: ControllerComponents) extends BaseController {

  private val cartList = mutable.ListBuffer(
    Cart(1, List(CartItem(1, 2), CartItem(2, 3))),
    Cart(2, List(CartItem(2, 2), CartItem(3, 1))),
    Cart(3, List(CartItem(1, 1), CartItem(3, 3)))
  )

  def getAll: Action[AnyContent] = Action {
    Ok(Json.toJson(cartList))
  }

  def getById(id: Int): Action[AnyContent] = Action {
    cartList.find(_.id == id) match {
      case Some(cart) => Ok(Json.toJson(cart))
      case None => NotFound(Json.obj("error" -> s"Cart with ID $id not found"))
    }
  }

  def add: Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[List[CartDTO]].fold(
      _ => BadRequest(Json.obj("error" -> "Incorrect data")),
      cartDTOList => {
        val newId = if (cartList.isEmpty) 1 else cartList.map(_.id).max + 1
        val cartWithId = Cart(newId, cartDTOList.map(dto => CartItem(dto.productId, dto.quantity)))
        cartList += cartWithId
        Created(Json.toJson(cartWithId))
      }
    )
  }

  def update(id: Int): Action[JsValue] = Action(parse.json) { request =>
    request.body.validate[List[CartDTO]].fold(
      _ => BadRequest(Json.obj("error" -> "Incorrect data")),
      updatedCartDTOList => {
        cartList.indexWhere(_.id == id) match {
          case -1 => NotFound(Json.obj("error" -> "Cart not found"))
          case idx =>
            val updatedCart = Cart(id, updatedCartDTOList.map(dto => CartItem(dto.productId, dto.quantity)))
            cartList.update(idx, updatedCart)
            Ok(Json.toJson(updatedCart))
        }
      }
    )
  }

  def delete(id: Int): Action[AnyContent] = Action {
    cartList.indexWhere(_.id == id) match {
      case -1 => NotFound(Json.obj("error" -> "Cart not found"))
      case idx =>
        cartList.remove(idx)
        NoContent
    }
  }
}
