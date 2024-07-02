/// <reference types="cypress" />

var spy;

Cypress.on('window:before:load', (window) => {
  spy = cy.spy(window.console, 'log');
});

describe('initialize the encrypted tunnel', () => {
    beforeEach(() => {
        cy.visit('http://localhost:5173')
    })
    
    it('checks that tunnel initialization message is logged', () => {
        // wait for the tunnel to initialize
        cy.wait(500)
        // check the console for the tunnel initialization message
        cy.then(() => {
            expect(spy).to.be.calledWith('[http://localhost:6001] Encrypted tunnel successfully established.')
        })
    })

    it('ensures that the page no longer contains a loading text', () => {
        // at the start of the test, the loading text should be present
        cy.get('.loader').should('exist')
        // wait for the tunnel to initialize
        cy.wait(500)
        // check that the loading text is no longer present
        cy.get('.loader').should('not.exist')
    })
})
