import 'mocha';
import {request, spec} from 'pactum';
import {faker} from '@faker-js/faker';
import {AuthUser, CreateDiscount, CreateUser, Discount} from "../gen";

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

    it('create user with same email', async () => {
        await spec()
            .post(`${baseUrl}/users`)
            .withJson(firstUser)
            .expectStatus(400)
    });

    it('should get all users', async () => {
        await spec()
            .get(`${baseUrl}/users`)
            .expectStatus(200)
            .expectJsonLength(2)
    });

    const authReq: AuthUser = {
        email: firstUser.email,
        password: firstUser.password
    }

    it('Login user', async () => {
        await spec()
            .post(`${baseUrl}/auth/token`)
            .withJson(authReq)
            .expectStatus(200)
            .expectJsonMatch({
                user: {
                    full_name: firstUser.full_name,
                    email: firstUser.email,
                    id: '$S{firstUserId}'
                }
            })
            .expectJsonSchema({
                type: 'object',
                required: ['token'],
            });
    });

    it('should login user with cookie', async () => {
        await spec()
            .post(`${baseUrl}/auth`)
            .withJson(authReq)
            .expectStatus(200)
            .expectCookiesLike({
                'Path': '/',
            })
            .expectJsonMatch({
                user: {
                    full_name: firstUser.full_name,
                    email: firstUser.email,
                    id: '$S{firstUserId}'
                }
            })
            .expectJsonSchema({
                type: 'object',
                required: ['token'],
            });
    })

    const discountFixed: CreateDiscount = {
        code: faker.string.alphanumeric(5),
        type: 'fixed',
        value: faker.number.int({min: 100, max: 1000}),
    }

    const discountPercentage: CreateDiscount = {
        code: faker.string.alphanumeric(5),
        type: 'percentage',
        value: faker.number.int({min: 1, max: 100}),
    }

    const discountPercentageInvalid: CreateDiscount = {
        code: faker.string.alphanumeric(5),
        type: 'percentage',
        value: faker.number.int({min: 101, max: 200}),
    }

    const discountWithUsageLimit: CreateDiscount = {
        code: faker.string.alphanumeric(5),
        type: 'fixed',
        value: faker.number.int({min: 100, max: 1000}),
        usage_limit: faker.number.int({min: 1, max: 10}),
    }

    const discountWithUsageLimitInvalid: CreateDiscount = {
        code: faker.string.alphanumeric(5),
        type: 'fixed',
        value: faker.number.int({min: 100, max: 1000}),
        usage_limit: faker.number.int({min: -10, max: -1}),
    }

    const discountWithExpiry: CreateDiscount = {
        code: faker.string.alphanumeric(5),
        type: 'fixed',
        value: faker.number.int({min: 100, max: 1000}),
        ends_at: faker.date.future().toISOString(),
        starts_at: faker.date.past().toISOString(),
    }

    const casesOK: { discount: CreateDiscount, message: string, expect: Discount }[] = [
        {
            discount: discountFixed,
            message: 'fixed discount',
            expect: {
                usage_limit: 0,
                usage_count: 0,
                ends_at: null,
                code: discountFixed.code,
                type: discountFixed.type,
                value: discountFixed.value,
                is_active: true
            }
        },
        {
            discount: discountPercentage,
            message: 'percentage discount',
            expect: {
                usage_limit: 0,
                usage_count: 0,
                ends_at: null,
                code: discountPercentage.code,
                type: discountPercentage.type,
                value: discountPercentage.value,
                is_active: true
            }
        },
        {
            discount: discountWithUsageLimit,
            message: 'discount with usage limit',
            expect: {
                usage_limit: discountWithUsageLimit.usage_limit,
                usage_count: 0,
                ends_at: null,
                code: discountWithUsageLimit.code,
                type: discountWithUsageLimit.type,
                value: discountWithUsageLimit.value,
                is_active: true
            }
        },
        {
            discount: discountWithExpiry,
            message: 'discount with expiry',
            expect: {
                usage_limit: 0,
                usage_count: 0,
                // @ts-ignore
                ends_at: discountWithExpiry.ends_at,
                code: discountWithExpiry.code,
                type: discountWithExpiry.type,
                value: discountWithExpiry.value,
                is_active: true
            }
        }
    ]

    const casesFail = [
        {
            discount: discountPercentageInvalid,
            message: 'percentage discount invalid',
        },
        {
            discount: discountWithUsageLimitInvalid,
            message: 'discount with usage limit invalid',
        },
    ]

    for (let c of casesOK) {
        it(`should create a ${c.message}`, async () => {
            await spec()
                .post(`${baseUrl}/discounts`)
                .withJson(c.discount)
                .expectStatus(201)
                .expectJsonMatch(c.expect)
                .expectJsonSchema({
                    type: 'object',
                    required: ['id', 'code', 'type', 'value', 'usage_limit', 'usage_count', 'created_at', 'updated_at'],
                });
        });
    }

    for (let c of casesFail) {
        it(`should fail to create a ${c.message}`, async () => {
            await spec()
                .post(`${baseUrl}/discounts`)
                .withJson(c.discount)
                .expectStatus(400)
        });
    }

    it('should get all discounts', async () => {
        await spec()
            .get(`${baseUrl}/discounts`)
            .expectStatus(200)
            .expectJsonLength(4)
    });
});