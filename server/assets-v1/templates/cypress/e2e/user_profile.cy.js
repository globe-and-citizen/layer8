describe('Authentication Page', () => {
    beforeEach(() => {
      cy.visit('http://localhost:5001/user')
    })
  
    it('loads the page correctly', () => {
      cy.get('title').should('contain', 'Authentication Page')

      cy.get('#app').should('exist')
    })
  
    it('verifies email when clicked', () => {
      cy.get('button').contains('Verify Email').click()
      cy.on('window:alert', (message) => {
        expect(message).to.equal('Email verified!')
      })
    })
  
    it('changes display name when clicked', () => {
      cy.get('input').contains('Display Name').type('New Name')
      cy.get('button').contains('Save change').click()
      cy.on('window:alert', (message) => {
        expect(message).to.equal('Display name changed!')
      })
    })
  
    it('logs out user when clicked', () => {
      cy.get('button').contains('Log out').click()
      cy.location('pathname').should('eq', '/') // Replace '/' with the expected URL after logout
    })
  })
  