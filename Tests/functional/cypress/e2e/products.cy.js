describe("Product Page Tests", () => {
  beforeEach(() => {
    cy.visit("/products");
  });

  it("The product list should be visible", () => {
    cy.get("ul").should("be.visible");
  });

  it('Should load the product list', () => {
    cy.get('h2').should('contain', 'Lista produktów')
    cy.get('ul li').should('have.length', 4)
    cy.get('ul li button')
        .filter(':contains("Dodaj do koszyka")')
        .should('have.length', 4)
  })

  it("each product list item should contain a dash '-' between name and price", () => {
    cy.get("ul li").each(($el) => {
      cy.wrap($el)
          .invoke("text")
          .should("match", /.+\s-\s.+/);
    });
  });

  it("Adding a product to the cart should update cart contents", () => {
    cy.get("button").contains("➕ Dodaj do koszyka").first().click();
    cy.get(".nav-link").contains("Koszyk").click();
    cy.url().should('include', '/cart')
    cy.get("ul").should("contain", "Ilość");
  });

  it("Product price should be displayed in Polish currency", () => {
    cy.get("li").first().should("contain", "zł");
  });

  it("Each product should have a name and price displayed", () => {
    cy.get("li").each(($el) => {
      cy.wrap($el).should("contain", "zł");
      cy.wrap($el).should("not.be.empty");
    });
  });

  it("User should be able to add multiple different products to the cart", () => {
    cy.visit("/products");

    cy.get("li").first().within(() => {
      cy.get("button").contains("➕ Dodaj do koszyka").click();
    });

    cy.get("li").eq(1).within(() => {
      cy.get("button").contains("➕ Dodaj do koszyka").click();
    });

    cy.get(".nav-link").contains("Koszyk").click();
    cy.url().should('include', '/cart')
    cy.wait(1000);

    cy.get("ul li").should("have.length.greaterThan", 1);
  });

  it("Clicking 'Add to Cart' multiple times should increase quantity", () => {
    cy.visit("/products");

    cy.get("button").contains("➕ Dodaj do koszyka").first().click();
    cy.wait(200);
    cy.get("button").contains("➕ Dodaj do koszyka").first().click();
    cy.wait(200);

    cy.get(".nav-link").contains("Koszyk").click();
    cy.get("ul").should("contain", "Ilość: 2");
  });


  it("Product names should should have add button", () => {
    cy.get("li").each(($el) => {
      cy.wrap($el).find("button").should("contain", "Dodaj do koszyka");
    });
  });

  it("should ensure user can navigate back to products from the cart", () => {
    cy.get(".nav-link").contains("Koszyk").click();
    cy.url().should("include", "/cart");
    cy.get(".nav-link").contains("Produkty").click();
    cy.url().should("include", "/products");
  });

});
