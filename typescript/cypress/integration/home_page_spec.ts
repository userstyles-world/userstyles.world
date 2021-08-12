
describe('The Home Page', () => {
    it('successfully loads', () => {
        cy.visit('/');
    });

    it('should have 5 elements in menu', () => {
        cy.get('.menu').children().should('have.length', 5);
    });

    it('should have the correct order in menu', () => {
        // The <li> elements are ordered in 'Search | Explore | Mod Log | Register | Sign in'
        cy.get('.menu').children().should(($menu) => {
            // Loop over the <li>'s children elements
            // And put the text in an array
            const menuText = [];
            $menu.children().each((i, el) => {
                if (0 === i) {
                    menuText.push((el as HTMLFormElement).className);
                } else {
                    menuText.push((el as HTMLAnchorElement).textContent);
                }
            });
            expect(menuText).to.deep.equal(['search', 'Explore', 'Mod Log', 'Register', 'Sign in']);
        });
    });

    it('should have the correct links in menu', () => {
        cy.get('.menu').children().should(($menu) => {
            // Loop over the <li>'s children elements
            // And put the text in an array
            const menuText = [];
            $menu.children().each((i, el) => {
                if (0 === i) {
                    menuText.push((el as HTMLFormElement).action);
                } else {
                    menuText.push((el as HTMLAnchorElement).href);
                }
            });
            expect(menuText).to.deep.equal([
                'http://localhost:3000/search',
                'http://localhost:3000/explore',
                'http://localhost:3000/modlog',
                'http://localhost:3000/register',
                'http://localhost:3000/login'
            ]);
        });
    });

    it('should have a explore button', () => {
        cy.get('.btn[href="/explore"]').should('have.length', 1);
    });
});

describe('Footer page', () => {
    it('should have the correct text', () => {
        cy.get('footer').should('exist');
    });
    it('should have an about section', () => {
        cy.get('footer').get('.about').contains('span', 'A free and open-source, community-driven website for browsing and sharing UserCSS userstyles.');
    });
    it('should have 2 more quarter section', () => {
        cy.get('footer').get('.Footer-list.quarter').should('have.length', 2);
    });
});


describe('The Home Page (DB seeding)', () => {
    // Ensure fresh test data.
    before(() => {
        cy.exec('cd ../ && ./tools/run testdata');
        // Ensure we have the test data.
        cy.visit('/');
    });

    it('should have 4 cards', () => {
        cy.get('.card').should('have.length', 4);
    });

    it('each card should have a picture element with images', () => {
        cy.get('.card').contains('a', 'picture');
    });

    it('should match that the stats has correct format', () => {
        cy.get('.stats').contains('UserStyles.world is home to 3 users and 6 userstyles. Styles have been viewed 0 times, installed 0 times, and updated 0 times in the last 24 hours, for a total of 0 views and 0 installs.');
    });
});
