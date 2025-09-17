'use client';

import { useState, useRef } from 'react';
import dynamic from 'next/dynamic';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ImageUpload } from '@/components/ui/image-upload';
import { MarkdownRenderer } from '@/components/ui/markdown-renderer';
import { Eye, Code, Image as ImageIcon, Link } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

// Dynamically import the markdown editor to avoid SSR issues
const MDEditor = dynamic(
  () => import('@uiw/react-md-editor'),
  { ssr: false }
);

interface MarkdownEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  height?: number;
  className?: string;
}

export function MarkdownEditor({ 
  value, 
  onChange, 
  placeholder = "Write your content here...",
  height = 400,
  className = ""
}: MarkdownEditorProps) {
  const [activeTab, setActiveTab] = useState<'edit' | 'preview'>('edit');
  const [linkUrl, setLinkUrl] = useState('');
  const [linkText, setLinkText] = useState('');
  const [showImageDialog, setShowImageDialog] = useState(false);
  const [showLinkDialog, setShowLinkDialog] = useState(false);
  const editorRef = useRef<any>(null);

  const insertImage = (imageUrl: string, altText?: string) => {
    const imageMarkdown = `![${altText || 'Image'}](${imageUrl})`;
    const newValue = value + (value ? '\n\n' : '') + imageMarkdown;
    onChange(newValue);
  };

  const insertLink = () => {
    if (linkUrl.trim() && linkText.trim()) {
      const linkMarkdown = `[${linkText}](${linkUrl})`;
      const newValue = value + (value ? '\n\n' : '') + linkMarkdown;
      onChange(newValue);
      setLinkUrl('');
      setLinkText('');
      setShowLinkDialog(false);
    }
  };


  return (
    <div className={`space-y-4 ${className}`}>
      {/* Custom Toolbar */}
      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <CardTitle className="text-sm font-medium">Rich Text Editor</CardTitle>
            <div className="flex items-center gap-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setShowImageDialog(true)}
                title="Insert Image"
              >
                <ImageIcon className="h-4 w-4" />
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setShowLinkDialog(true)}
                title="Insert Link"
              >
                <Link className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-0">
          <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'edit' | 'preview')}>
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="edit" className="flex items-center gap-2">
                <Code className="h-4 w-4" />
                Edit
              </TabsTrigger>
              <TabsTrigger value="preview" className="flex items-center gap-2">
                <Eye className="h-4 w-4" />
                Preview
              </TabsTrigger>
            </TabsList>
            
            <TabsContent value="edit" className="mt-4">
              <div className="border rounded-md">
                <MDEditor
                  ref={editorRef}
                  value={value}
                  onChange={(val) => onChange(val || '')}
                  height={height}
                  data-color-mode="light"
                  preview="edit"
                />
              </div>
            </TabsContent>
            
            <TabsContent value="preview" className="mt-4">
              <div className="border rounded-md p-4 min-h-[400px]">
                <MarkdownRenderer content={value} />
              </div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* Image Upload Dialog */}
      {showImageDialog && (
        <ImageUpload
          onImageInsert={insertImage}
          onClose={() => setShowImageDialog(false)}
        />
      )}

      {/* Link Insert Dialog */}
      {showLinkDialog && (
        <Card className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <Card className="w-96">
            <CardHeader>
              <CardTitle>Insert Link</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label htmlFor="linkText">Link Text</Label>
                <Input
                  id="linkText"
                  value={linkText}
                  onChange={(e) => setLinkText(e.target.value)}
                  placeholder="Click here"
                />
              </div>
              <div>
                <Label htmlFor="linkUrl">URL</Label>
                <Input
                  id="linkUrl"
                  value={linkUrl}
                  onChange={(e) => setLinkUrl(e.target.value)}
                  placeholder="https://example.com"
                />
              </div>
              <div className="flex gap-2">
                <Button onClick={insertLink} disabled={!linkUrl.trim() || !linkText.trim()}>
                  Insert Link
                </Button>
                <Button variant="outline" onClick={() => setShowLinkDialog(false)}>
                  Cancel
                </Button>
              </div>
            </CardContent>
          </Card>
        </Card>
      )}
    </div>
  );
}
