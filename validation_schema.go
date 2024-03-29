package main

var schema = `
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
	  "order_uid": {"type": "string"},
	  "track_number": {"type": "string"},
	  "entry": {"type": "string"},
	  "delivery": {
		"type": "object",
		"properties": {
		  "name": {"type": "string"},
		  "phone": {"type": "string"},
		  "zip": {"type": "string"},
		  "city": {"type": "string"},
		  "address": {"type": "string"},
		  "region": {"type": "string"},
		  "email": {"type": "string", "format": "email"}
		},
		"required": ["name", "phone", "zip", "city", "address", "region", "email"]
	  },
	  "payment": {
		"type": "object",
		"properties": {
		  "transaction": {"type": "string"},
		  "request_id": {"type": "string"},
		  "currency": {"type": "string"},
		  "provider": {"type": "string"},
		  "amount": {"type": "integer"},
		  "payment_dt": {"type": "integer"},
		  "bank": {"type": "string"},
		  "delivery_cost": {"type": "integer"},
		  "goods_total": {"type": "integer"},
		  "custom_fee": {"type": "integer"}
		},
		"required": ["transaction", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee"]
	  },
	  "items": {
		"type": "array",
		"items": {
		  "type": "object",
		  "properties": {
			"chrt_id": {"type": "integer"},
			"track_number": {"type": "string"},
			"price": {"type": "integer"},
			"rid": {"type": "string"},
			"name": {"type": "string"},
			"sale": {"type": "integer"},
			"size": {"type": "string"},
			"total_price": {"type": "integer"},
			"nm_id": {"type": "integer"},
			"brand": {"type": "string"},
			"status": {"type": "integer"}
		  },
		  "required": ["chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status"]
		}
	  },
	  "locale": {"type": "string"},
	  "internal_signature": {"type": "string"},
	  "customer_id": {"type": "string"},
	  "delivery_service": {"type": "string"},
	  "shardkey": {"type": "string"},
	  "sm_id": {"type": "integer"},
	  "date_created": {"type": "string", "format": "date-time"},
	  "oof_shard": {"type": "string"}
	},
	"required": ["order_uid", "track_number", "entry", "delivery", "payment", "items", "locale", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"]
  }
  `
