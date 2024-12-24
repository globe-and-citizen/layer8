/// <reference types="cypress" />

describe('authenticate user to WGP', () => {
    beforeEach(() => {
        cy.visit('http://localhost:5173')
    })

    it('registers a new user', () => {
        // clicking the "don't have an account" link
        cy.get('a[class="block"]').click()
        
        // filling in the registration form
        cy.get('input').eq(0).type('testuser')
        cy.get('input').eq(1).type('testpassword')
        // uploading a profile picture
        cy.get('input[type="file"]').selectFile('cypress/fixtures/profile.jpg', { force: true })

        // wait for the element with "loading" class to disappear
        cy.wait(500)

        // submitting the form
        cy.get('button').eq(1).click()

        // wait for the registration to complete
        cy.wait(500)

        // check registration success alert
        const stub = cy.stub()
        cy.on('window:alert', stub)
        cy.then(() => {
            expect(stub).to.be.calledWith('Registration successful!')
        })
    })

    it('logs in the registered user', () => {
        // filling in the login form
        cy.get('input').eq(0).type('testuser')
        cy.get('input').eq(1).type('testpassword')

        // wait for the element with "loading" class to disappear
        cy.wait(500)

        // submitting the form
        cy.get('button').eq(1).click()

        // wait for the login to complete
        cy.wait(500)

        // check login success alert
        const stub = cy.stub()
        cy.on('window:alert', stub)
        cy.then(() => {
            expect(stub).to.be.calledWith('Login successful!')
        })
    })
    
    it('logs in anonymously', () => {
        cy.loginWGP()
        
        cy.contains('Login Anonymously').click()

        // check anonymous login success alert
        const stub = cy.stub()
        cy.on('window:alert', stub)
        cy.then(() => {
            expect(stub).to.be.calledWith('You are now logged in anonymously!')
        })

        // wait for the welcome page to load
        cy.wait(500)

        // check that the welcome page is displayed
        cy.get('h1').should('contain', 'Welcome testuser!')
    })

    it('logs in with oauth', () => {
        cy.loginWGP()

        cy.contains('Login with Layer8 Redirect').click()
        cy.get('input[name="username"]').type("tester")
        cy.get('input[name="password"]').type("12341234")
        cy.get('input[type="submit"]').click()
        
        // authorize the app
        cy.contains("Authorize").should("exist")
        cy.get('input[type="submit"]').click()

        // ensure the app is redirected to the callback page
        cy.url().should('contain', '/oauth2/callback')

        // wait for the welcome page to load
        cy.wait(500)

        // check that the welcome page is displayed
        cy.get('h1').should('contain', 'Welcome testuser!')
    })

    it('logs out the user', () => {
        cy.loginAnonymouslyWGP()

        // logging out the user
        cy.contains('Logout').click()

        // check that the login page is displayed
        cy.get('h2').should('contain', 'Login')
    })
})