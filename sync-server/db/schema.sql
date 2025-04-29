CREATE TABLE "users" (
                         "id" serial PRIMARY KEY NOT NULL,
                         "name" text NOT NULL,
                         "created_at" timestamp DEFAULT now()
);
CREATE TABLE "session_to_user_mapping" (
                                           "id" serial PRIMARY KEY NOT NULL,
                                           "session_id" uuid NOT NULL,
                                           "user_id" bigint NOT NULL
);
--> statement-breakpoint
CREATE TABLE "sessions" (
                            "id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
                            "content" "bytea",
                            "whiteboard_content" "bytea"
);
--> statement-breakpoint
ALTER TABLE "users" ADD COLUMN "username" text NOT NULL;--> statement-breakpoint
ALTER TABLE "users" ADD COLUMN "email" text NOT NULL;--> statement-breakpoint
ALTER TABLE "users" ADD COLUMN "password" text;--> statement-breakpoint
ALTER TABLE "session_to_user_mapping" ADD CONSTRAINT "session_to_user_mapping_session_id_sessions_id_fk" FOREIGN KEY ("session_id") REFERENCES "public"."sessions"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session_to_user_mapping" ADD CONSTRAINT "session_to_user_mapping_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "users" ADD CONSTRAINT "users_username_unique" UNIQUE("username");--> statement-breakpoint
ALTER TABLE "users" ADD CONSTRAINT "users_email_unique" UNIQUE("email");
ALTER TABLE "session_to_user_mapping" ADD COLUMN "is_admin" boolean NOT NULL;
ALTER TABLE "session_to_user_mapping" DROP CONSTRAINT "session_to_user_mapping_user_id_users_id_fk";
--> statement-breakpoint
ALTER TABLE "session_to_user_mapping" ADD COLUMN "user_email" text NOT NULL;--> statement-breakpoint
ALTER TABLE "sessions" ADD COLUMN "name" text NOT NULL;--> statement-breakpoint
ALTER TABLE "sessions" ADD COLUMN "is_anyone_allowed" boolean NOT NULL;--> statement-breakpoint
ALTER TABLE "session_to_user_mapping" ADD CONSTRAINT "session_to_user_mapping_user_email_users_email_fk" FOREIGN KEY ("user_email") REFERENCES "public"."users"("email") ON DELETE no action ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session_to_user_mapping" DROP COLUMN "user_id";