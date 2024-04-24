import { AuthUser } from "./AuthUser";
import type { UserWithToken } from "./UserWithToken";

 /**
 * @description OK
*/
export type PostAdminAuth200 = UserWithToken;

 /**
 * @description User data
*/
export type PostAdminAuthMutationRequest = AuthUser;

 /**
 * @description OK
*/
export type PostAdminAuthMutationResponse = UserWithToken;

 export type PostAdminAuthMutation = {
    Response: PostAdminAuthMutationResponse;
    Request: PostAdminAuthMutationRequest;
};