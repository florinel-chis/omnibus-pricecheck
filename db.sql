CREATE TABLE `prices` (
  `price_id` int NOT NULL AUTO_INCREMENT,
  `sku` varchar(32) NOT NULL,
  `date` datetime NOT NULL,
  `list_price` decimal(10,4) NOT NULL,
  `final_price` decimal(10,4) NOT NULL,
  PRIMARY KEY (`price_id`),
  UNIQUE KEY `idx_sku_date` (`sku`,`date`),
  KEY `sku` (`sku`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci 