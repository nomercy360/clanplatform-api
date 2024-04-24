import 'mocha';
import {request, spec} from 'pactum';
import {faker} from '@faker-js/faker';
import {CreateUser} from "../gen";

const baseUrl = "http://localhost:8080/admin";

describe('Test Admin Routes', () => {
    before(async () => {
        request.setDefaultTimeout(10000);
    });

    const firstUser: CreateUser = {
        full_name: faker.person.fullName(),
        email: faker.internet.email(),
        password: faker.internet.password(),
    };

    const secondUser: CreateUser = {
        full_name: faker.person.fullName(),
        email: faker.internet.email(),
        password: faker.internet.password()
    }

    for (let user of [firstUser, secondUser]) {
        it('should create a user', async () => {
            let storeVarName = user === firstUser ? 'firstUserId' : 'secondUserId';
            await spec()
                .post(`${baseUrl}/users`)
                .withJson(user)
                .expectStatus(201)
                .expectJsonMatch({
                    full_name: user.full_name,
                    email: user.email,
                })
                .expectJsonSchema({
                    type: 'object',
                    required: ['id', 'full_name', 'email', 'created_at', 'updated_at'],
                })
                .stores(storeVarName, 'id');
        });
    }

    it('should get all users', async () => {
        await spec()
            .get(`${baseUrl}/users`)
            .expectStatus(200)
            .expectJsonLength(2)
    });

    it('Login user', async () => {
        await spec()
            .post(`${baseUrl}/auth/token`)
            .withJson({
                "email": firstUser.email,
                "password": firstUser.password
            })
            .expectStatus(200)
            .expectJsonSchema({
                type: 'object',
                required: ['token'],
            });
    });
});