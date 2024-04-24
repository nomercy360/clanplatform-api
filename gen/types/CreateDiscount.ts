export const createDiscountType = {
    "percentage": "percentage",
    "fixed": "fixed",
    "free_shipping": "free_shipping"
} as const;
export type CreateDiscountType = (typeof createDiscountType)[keyof typeof createDiscountType];
export type CreateDiscount = {
    /**
     * @type string
    */
    code: string;
    /**
     * @type string | undefined
    */
    ends_at?: string;
    /**
     * @type string | undefined
    */
    starts_at?: string;
    /**
     * @type string
    */
    type: CreateDiscountType;
    /**
     * @type integer | undefined
    */
    usage_limit?: number;
    /**
     * @type integer
    */
    value: number;
};