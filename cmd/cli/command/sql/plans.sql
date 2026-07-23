INSERT INTO plans (
    id,
    slug,
    name,
    description,
    billing_interval,
    price,
    currency,
    stripe_price_id,
    is_active,
    created_at,
    updated_at
)
VALUES
('e75e7940-2683-4ea1-a5c7-6d1724892369','free','Free','Découverte du produit, analyses limitées','monthly',0,'eur','',true,now(),now()),
('e437e948-9ac7-422a-8a86-79ca43ca9d5d','starter','Starter','Pour un usage régulier, image uniquement','monthly',1200,'eur','price_starter_monthly',true,now(),now()),
('aca09dc8-1c60-4670-9e2a-2eac0a8dc268','starter','Starter','Pour un usage régulier, image uniquement','annually',11500,'eur','price_starter_annually',true,now(),now()),
('dcad50c5-7a2d-4de6-958f-e9d7827fab4d','pro','Pro','Pipeline complet, images et vidéos','monthly',3900,'eur','price_pro_monthly',true,now(),now()),
('33c61c6e-d2c5-4915-9b5a-3fa8287047b8','pro','Pro','Pipeline complet, images et vidéos','annually',37500,'eur','price_pro_annually',true,now(),now()),
('b8f650d8-e625-4604-ac34-033de9f5246a','business','Business','Volume, support dédié, détail par sous-signal','monthly',14900,'eur','price_business_monthly',true,now(),now()),
('c90170c4-3824-4190-b81b-207db611442b','business','Business','Volume, support dédié, détail par sous-signal','annually',143000,'eur','price_business_annually',true,now(),now())
ON CONFLICT (id) DO NOTHING;
