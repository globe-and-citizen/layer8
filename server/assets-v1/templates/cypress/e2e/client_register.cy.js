describe('Register Client Page', () => {
    beforeEach(() => {
      cy.visit('http://localhost:5001/client-register-page')
    })
  
    it('displays the registration form', () => {
      cy.get('h1').should('contain', 'Register your product')
      cy.get('input[id="name"]').should('exist')
      cy.get('input[id="redirect_uri"]').should('exist')
      cy.get('input[id="username"]').should('exist')
      cy.get('input[id="password"]').should('exist')
      cy.get('button').should('contain', 'Register')
      cy.contains('Already have an account?').should('exist')
    })
  
    it('allows clients to register with valid data', () => {
      cy.get('input[id="name"]').type('Test Project')
      cy.get('input[id="redirect_uri"]').type('https://example.com/callback')
      cy.get('input[id="username"]').type('testuser')
      cy.get('input[id="password"]').type('password123')
      cy.get('button').click()
      cy.url().should('include', 'http://localhost:5001/client-login-page')
    })
  
    it('displays an error message for incomplete registration data', () => {
      cy.get('button').click()
      cy.on('window:alert', (message) => {
        expect(message).to.equal('Please enter all fields!')
      })
    })
  
    // it('displays an error message for invalid project name', () => {
    //   cy.get('input[id="name"]').type(' ') // Empty project name
    //   cy.get('input[id="redirect_uri"]').type('https://example.com/callback')
    //   cy.get('input[id="username"]').type('testuser')
    //   cy.get('input[id="password"]').type('password123')
    //   cy.get('button').click()
    //   cy.on('window:alert', (message) => {
    //     expect(message).to.equal('Please enter all fields!')
    //   })
    // })
  
    // it('displays an error message for invalid redirect URL', () => {
    //     cy.get('input[id="name"]').type('Test Project')
    //     cy.get('input[id="redirect_uri"]').type('invalid-url') // Invalid URL format
    //     cy.get('input[id="username"]').type('testuser')
    //     cy.get('input[id="password"]').type('password123')
    //     cy.get('button').click()
    //     cy.on('window:alert', (message) => {
    //       expect(message).to.equal('Please enter all fields!')
    //     })
    //   })
      
    //   it('displays an error message for invalid username', () => {
    //     cy.get('input[id="name"]').type('Test Project')
    //     cy.get('input[id="redirect_uri"]').type('https://example.com/callback')
    //     cy.get('input[id="username"]').type('') // Empty username
    //     cy.get('input[id="password"]').type('password123')
    //     cy.get('button').click()
    //     cy.on('window:alert', (message) => {
    //       expect(message).to.equal('Please enter all fields!')
    //     })
    //   })
      
    //   it('displays an error message for invalid password', () => {
    //     cy.get('input[id="name"]').type('Test Project')
    //     cy.get('input[id="redirect_uri"]').type('https://example.com/callback')
    //     cy.get('input[id="username"]').type('testuser')
    //     cy.get('input[id="password"]').type('') // Empty password
    //     cy.get('button').click()
    //     cy.on('window:alert', (message) => {
    //       expect(message).to.equal('Please enter all fields!')
    //     })
    //   })
      
  })
  