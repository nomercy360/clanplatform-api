import { CreateDiscount } from "./CreateDiscount";
import type { Discount } from "./Discount";

 /**
 * @description OK
*/
export type PostAdminDiscounts200 = Discount;

 /**
 * @description Discount data
*/
export type PostAdminDiscountsMutationRequest = CreateDiscount;

 /**
 * @description OK
*/
export type PostAdminDiscountsMutationResponse = Discount;

 export type PostAdminDiscountsMutation = {
    Response: PostAdminDiscountsMutationResponse;
    Request: PostAdminDiscountsMutationRequest;
};