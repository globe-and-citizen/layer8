/// <reference types="cypress" />

var spy;

Cypress.on('window:before:load', (window) => {
  spy = cy.spy(window.console, 'log');
});

describe('initialize the encrypted tunnel', () => {
    it('checks that tunnel initialization message is logged', () => {
        cy.visit('http://localhost:5173')
        // wait for the tunnel to initialize
        cy.wait(500)
        // check the console for the tunnel initialization message
        cy.then(() => {
            expect(spy).to.be.calledWith('[http://localhost:6002] Encrypted tunnel successfully established.')
        })
    })
})
