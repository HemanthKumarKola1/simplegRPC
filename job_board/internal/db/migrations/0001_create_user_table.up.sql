CREATE TABLE public.user (
    userID SERIAL PRIMARY KEY,
    userName TEXT UNIQUE NOT NULL,
    passwordHash TEXT NOT NULL,
    email TEXT NOT NULL,
    userType INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_type_check CHECK (UserType IN (1, 2))
);