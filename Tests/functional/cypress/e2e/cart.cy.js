describe("Shopping Cart Tests", () => {
    beforeEach(() => {
        cy.visit("/cart");
    });


    it("The cart should be empty initially", () => {
        cy.get("p").should("contain", "Koszyk jest pusty");
    });


    it("Adding a product should increase quantity in the cart", () => {
        cy.visit("/products");
        cy.get("button").contains("➕ Dodaj do koszyka").first().click();
        cy.get(".nav-link").contains("Koszyk").click();
        cy.get("ul li").should("have.length.greaterThan", 0);
    });

    it("Total cart value should be displayed in Polish", () => {
        cy.get("h3").should("contain", "Całkowita wartość koszyka");
    });

});
