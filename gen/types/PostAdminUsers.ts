import { CreateUser } from "./CreateUser";
import type { User } from "./User";

 /**
 * @description Created
*/
export type PostAdminUsers201 = User;

 /**
 * @description User data
*/
export type PostAdminUsersMutationRequest = CreateUser;

 /**
 * @description Created
*/
export type PostAdminUsersMutationResponse = User;

 export type PostAdminUsersMutation = {
    Response: PostAdminUsersMutationResponse;
    Request: PostAdminUsersMutationRequest;
};