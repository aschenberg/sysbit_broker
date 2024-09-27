ALTER TABLE lesson_page ADD CONSTRAINT lessonpage_lessonid_fkey FOREIGN KEY(lesson_id) 
      REFERENCES lesson(lesson_id) ON DELETE CASCADE;
ALTER TABLE lesson_content ADD CONSTRAINT lessoncontent_pageid_fkey FOREIGN KEY(page_id) 
      REFERENCES lesson_page(page_id) ON DELETE CASCADE;
