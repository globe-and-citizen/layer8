/// <reference types="cypress" />

describe('get next poem', () => {
    beforeEach(() => {
        cy.loginAnonymouslyWGP();
    })

    it('gets the next poem', () => {
        // the poem should be empty at the start
        cy.get('td').eq(0).should('have.text', '')
        cy.get('td').eq(1).should('have.text', '')
        cy.get('td').eq(2).should('have.text', '')

        // get the next poem
        cy.contains('Get Next Poem').click()
        cy.wait(50)

        cy.get('td').eq(0).then(($el) => {
            const poem = $el.text()
            expect(poem).to.not.equal('')
        })
        cy.get('td').eq(1).then(($el) => {
            const poem = $el.text()
            expect(poem).to.not.equal('')
        })
        cy.get('td').eq(2).then(($el) => {
            const poem = $el.text()
            expect(poem).to.not.equal('')
        })
    })
})
