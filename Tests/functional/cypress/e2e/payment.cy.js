describe("Payment Page Tests", () => {
    beforeEach(() => {
        cy.visit("/payment");
    });

    it("The payment form should be visible", () => {
        cy.get("form").should("be.visible");
    });

    it("User should be able to enter payment details and submit the form", () => {
        cy.get("input[placeholder='Imię']").type("Jan");
        cy.get("input[placeholder='Nazwisko']").type("Kowalski");
        cy.get("input[placeholder='Numer karty']").type("1234 5678 9012 3456");
        cy.get("button").contains("Zapłać").click();
        cy.on("window:alert", (txt) => {
            expect(txt).to.contains("Płatność przetworzona!");
        });
    });

    it("Payment amount should match the total cart value", () => {
        cy.visit("/products");
        cy.get("button").contains("➕ Dodaj do koszyka").first().click();

        cy.visit("/cart");
        cy.get("h3").invoke("text").then((cartTotal) => {
            const totalValue = cartTotal.match(/\d+(\.\d{1,2})?/)[0];
            cy.visit("/payment");
            cy.get("h3").should("contain", totalValue + " zł");
        });
    });

    it("Payment should not proceed when name and surname fields are empty", () => {
        cy.get("input[placeholder='Numer karty']").type("1234 5678 9012 3456");
        cy.get("button").contains("Zapłać").click();
        cy.on("window:alert", (txt) => {
            expect(txt).to.contains("Błąd płatności");
        });
    });

    it("Credit card number should not accept non-numeric values", () => {
        cy.get("input[placeholder='Numer karty']").type("ABCD XYZ");
        cy.get("button").contains("Zapłać").click();
        cy.on("window:alert", (txt) => {
            expect(txt).to.contains("Błąd płatności");
        });
    });


});
