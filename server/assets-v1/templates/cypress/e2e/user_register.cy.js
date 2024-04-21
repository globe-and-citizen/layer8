describe('Authentication Page', () => {
    beforeEach(() => {
      cy.visit('http://localhost:5001/user-register-page')
    })
  
    it('displays the registration form', () => {
      cy.get('h1').should('contain', 'Register')
      cy.get('input[placeholder="Email"]').should('exist')
      cy.get('input[placeholder="Username"]').should('exist')
      cy.get('input[placeholder="First Name"]').should('exist')
      cy.get('input[placeholder="Last Name"]').should('exist')
      cy.get('input[placeholder="Display Name"]').should('exist')
      cy.get('input[placeholder="Country"]').should('exist')
      cy.get('input[type="password"]').should('exist')
      cy.get('button').should('contain', 'Register')
    })
  
    it('allows users to register with valid data', () => {
      cy.get('input[placeholder="Email"]').type('test@example.com')
      cy.get('input[placeholder="Username"]').type('testuser')
      cy.get('input[placeholder="First Name"]').type('John')
      cy.get('input[placeholder="Last Name"]').type('Doe')
      cy.get('input[placeholder="Display Name"]').type('JohnDoe')
      cy.get('input[placeholder="Country"]').type('United States')
      cy.get('input[type="password"]').type('password123')
      cy.get('button').click()
      cy.url().should('include', 'http://localhost:5001/user-login-page')
    })
  
    it('displays an error message for incomplete registration data', () => {
      cy.get('button').click()
      cy.on('window:alert', (message) => {
        expect(message).to.equal('Please enter a username and password!')
      })
    })
  })
  