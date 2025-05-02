describe("Payment Page Tests", () => {
    beforeEach(() => {
        cy.visit("/payment");
    });

    it("The payment form should be visible", () => {
        cy.get("form").should("be.visible");
    });

    it("should display payment header with correct text", () => {
        cy.get("h2")
            .should("be.visible")
            .should("have.text", "Formularz płatności")
            .should("contain", "płatności");
    });

    it("all input fields in the payment form should be visible and empty initially", () => {
        cy.get("input[placeholder='Imię']")
            .should("be.visible")
            .should("have.value", "");
        cy.get("input[placeholder='Nazwisko']")
            .should("be.visible")
            .should("have.value", "");
        cy.get("input[placeholder='Numer karty']")
            .should("be.visible")
            .should("have.value", "");
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
            expect(txt)
                .should("be.a", "string")
                .should("contain", "Błąd płatności");
        });
    });

    it("Credit card number should not accept non-numeric values", () => {
        cy.get("input[placeholder='Numer karty']").type("ABCD XYZ");
        cy.get("button").contains("Zapłać").click();
        cy.on("window:alert", (txt) => {
            expect(txt)
                .should("be.a", "string")
                .should("contain", "Błąd płatności");
        });
    });

    it("should reset the payment amount to 0 after successful payment", () => {

        cy.visit("/products");
        cy.get("button").contains("➕ Dodaj do koszyka").first().click();
        cy.visit("/payment");

        cy.get("input[placeholder='Imię']").type("Jan");
        cy.get("input[placeholder='Nazwisko']").type("Kowalski");
        cy.get("input[placeholder='Numer karty']").type("1234 5678 9012 3456");
        cy.get("button").contains("Zapłać").click();

        cy.get("h3").should("be.visible").and("contain", "0.00 zł").and("have.text", "Całkowita wartość płatności: 0.00 zł");
    });

    it("should clear the cart after successful payment", () => {
        cy.get("input[placeholder='Imię']")
            .type("Jan")
            .should("have.value", "Jan");

        cy.get("input[placeholder='Nazwisko']").type("Kowalski").should("have.value", "Kowalski");

        cy.get("input[placeholder='Numer karty']").type("1234 5678 9012 3456").should("have.value", "1234 5678 9012 3456");

        cy.get("button").contains("Zapłać").click();
        cy.visit("/cart");
        cy.get("p").should("be.visible")
            .and("contain", "Koszyk jest pusty!")
            .and("not.contain", "Ilość")
            .and("not.contain", "zł");
    });



});
