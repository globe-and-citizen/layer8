const project_name = 'Test Project';
const username = 'testuser';
const password = 'password123';
const redirect_uri = 'https://redirecturl.com';
const backend_uri = 'https://backendurl.com';

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
    cy.get('input[id="name"]').type(project_name)
    cy.get('input[id="redirect_uri"]').type(redirect_uri)
    cy.get('input[id="backend_uri"]').type(backend_uri)
    cy.get('input[id="username"]').type(username)
    cy.get('input[id="password"]').type(password)
    cy.get('button').click()
    cy.url().should('include', 'http://localhost:5001/client-login-page')
  })

  it('displays an error message for incomplete registration data', () => {
    cy.get('button').click()
    
    cy.get('.bg-red-500').should('be.visible')
    cy.contains('.bg-red-500', 'Please enter all fields!').should('be.visible')
  })
})

describe('Login Page', () => {
  beforeEach(() => {
    cy.visit('http://localhost:5001/client-login-page')
  })

  it('displays the login form', () => {
    cy.get('h1').should('contain', 'Login')
    cy.get('input#username').should('exist')
    cy.get('input#password').should('exist')
    cy.get('button').should('contain', 'Login')
  })

  it('allows users to login with valid credentials', () => {
    cy.get('input#username').type(username)
    cy.get('input#password').type(password)
    cy.get('button').click()
    cy.url().should('include', 'http://localhost:5001/client-profile')
    cy.contains('Welcome “Test Project!” Client Portal').should('be.visible');
    cy.contains('Your data').should('be.visible');               
    
    cy.contains('.font-bold', 'Name:').next(). should('exist');
    cy.get('input[placeholder="Name"]').should('have.value', project_name);
    
    cy.contains('.font-bold', 'Redirect URI:').next().should('exist');
    cy.get('input[placeholder="Redirect URI"]').should('have.value', redirect_uri);

    cy.contains('.font-bold', 'UUID:').next().should('exist');
    cy.get('input[placeholder="UUID"]').should('exist').should(($input) => {
      expect($input.val()).not.to.be.empty;
    });
    
    cy.contains('.font-bold', 'Secret:').should('exist');
    cy.get('input[placeholder="Secret"]').should('exist').should(($input) => {
      expect($input.val()).not.to.be.empty;
    });
  })

  it('Copying UUID to clipboard', () => {
    cy.get('input#username').type(username);
    cy.get('input#password').type(password);
    cy.get('button').click();
    cy.url().should('include', 'http://localhost:5001/client-profile');
  
    const userId = 'bd2422b6-2357-4f8f-ba46-c1e70c5f0173';
    
    cy.window().then((window) => {
      window.document.execCommand = cy.stub().returns(true);
      window.user = {
        id: userId
      };
    });
  
    cy.get('button[value="UUID"]').click();
  
    cy.wait(1000);
  
    cy.window().then((window) => {
      expect(window.document.execCommand).to.have.been.calledOnceWith('copy');
    });
  });

  it('Copying Secret to clipboard', () => {
    cy.get('input#username').type(username);
    cy.get('input#password').type(password);
    cy.get('button').click();
    cy.url().should('include', 'http://localhost:5001/client-profile');
  
    const secret = 'b333a024c425f1b250e9cd8084093220edbddc7f727ab31797232e48a3d57a59';
    
    cy.window().then((window) => {
      window.document.execCommand = cy.stub().returns(true);
      window.user = {
        id: secret
      };
    });
  
    cy.get('button[value="Secret"]').click();
  
    cy.wait(1000);
  
    cy.window().then((window) => {
      expect(window.document.execCommand).to.have.been.calledOnceWith('copy');
    });
  });
})

describe('Delete User', () => {
  const tableName = 'clients';

  it('deletes the registered client', () => {
    cy.deleteRegisteredUser(username, tableName).then((result) => {
      expect(result).to.be.true; 
    });
  });
});