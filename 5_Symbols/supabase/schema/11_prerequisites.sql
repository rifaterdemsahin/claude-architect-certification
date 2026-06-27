-- Create the prerequisites table
CREATE TABLE IF NOT EXISTS public.prerequisites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_number INTEGER NOT NULL,
    video_number INTEGER NOT NULL,
    install_name TEXT NOT NULL,
    install_command TEXT NOT NULL,
    verification_command TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now()) NOT NULL
);

-- Set up Row Level Security (RLS)
ALTER TABLE public.prerequisites ENABLE ROW LEVEL SECURITY;

-- Create policies for public access (or authenticated users depending on requirements)
-- Assuming we want this readable by anyone for now, like course outline
CREATE POLICY "Enable read access for all users" ON public.prerequisites FOR SELECT USING (true);
CREATE POLICY "Enable all access for authenticated users" ON public.prerequisites FOR ALL USING (auth.role() = 'authenticated');
