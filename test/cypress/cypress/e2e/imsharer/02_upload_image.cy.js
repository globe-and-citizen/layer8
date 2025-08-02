const images = Array.from({ length: 6 }, (_, i) => `image${i + 1}.jpg`)

describe('upload and display images', () => {
    beforeEach(() => {
        cy.visit('http://localhost:5173')
        // wait for the tunnel to initialize
        cy.wait(10000)
    })

    it('uploads images', () => {
        images.forEach((image) => {
            cy.get('input[type="file"]').selectFile(`cypress/fixtures/${image}`, { force: true })
            cy.get('input[type="button"]').click()
        })
    })

    it('displays the images in the gallery', () => {
        // the number of images in the gallery should be at least the number of images uploaded
        cy.get('.gallery').find('img').should('have.length.gte', images.length)
        // each image should have an alt attribute
        images.forEach((image) => {
            cy.get(`img[alt="${image}"]`).should('exist')
        })
    })
})