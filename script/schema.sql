CREATE TABLE `template` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `template_name` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `parse_type` int NOT NULL,
  `auto_gen_object_id` bool NOT NULL,
  `origin_object_id` varchar(255) NOT NULL,
  `value_extract` varchar(255) NOT NULL,
  `go_struct` varchar(2550) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `objectids` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `template_id` int NOT NULL,
  `serial` int NOT NULL,
  `object_id` varchar(255) UNIQUE NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `objectids` ADD FOREIGN KEY (`template_id`) REFERENCES `template` (`id`);
