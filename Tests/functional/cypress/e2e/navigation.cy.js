describe("Navigation Bar Tests", () => {
    beforeEach(() => {
        cy.visit("/");
    });

    it("Navigation bar should be visible", () => {
        cy.get(".navbar").should("be.visible");
    });

    it("Home button should work and display the correct heading", () => {
        cy.get(".nav-link").contains("Strona główna").should("be.visible").click();
        cy.url().should("eq", Cypress.config("baseUrl") + "/");
        cy.get("h1").should("contain", "Witaj w sklepie");
    });

    it("Products button should work and show product list", () => {
        cy.get(".nav-link").contains("Produkty").should("be.visible").click();
        cy.url().should("eq", Cypress.config("baseUrl") + "/products");
        cy.get("h2").should("contain", "Lista produktów");
        cy.get("ul").should("be.visible");
    });

    it("Cart button should work and indicate empty cart or items", () => {
        cy.get(".nav-link").contains("Koszyk").should("be.visible").click();
        cy.url().should("eq", Cypress.config("baseUrl") + "/cart");
        cy.get("h2").should("contain", "Twój koszyk");
        cy.get("p, h3").should("be.visible");
    });

    it("Payment button should work and display the payment form", () => {
        cy.get(".nav-link").contains("Płatności").should("be.visible").click();
        cy.url().should("eq", Cypress.config("baseUrl") + "/payment");
        cy.get("h2").should("contain", "Formularz płatności");
        cy.get("form").should("be.visible");
    });
});
