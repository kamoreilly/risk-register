-- Add responder role to user_role enum for incident management
ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'responder';
