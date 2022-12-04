package micropub

import (
	"fmt"
	"tiim/go-comment-api/lib/mfobjects"
	"time"

	"willnorris.com/go/microformats"
)

func ModifyEntry(post *MicropubPost, deleteProps interface{}, addProps, replaceProps map[string][]interface{}) error {

	mf := post.Entry.ToMicroformat()
	for key, values := range replaceProps {
		mf.Properties[key] = values
	}
	for key, values := range addProps {
		_, ok := mf.Properties[key]
		if !ok {
			mf.Properties[key] = values
		} else {
			mf.Properties[key] = append(mf.Properties[key], values...)
		}
	}
	if deleteProps != nil {
		switch del := deleteProps.(type) {
		case []interface{}:
			for _, key := range del {
				delete(mf.Properties, key.(string))
			}
		case map[string]interface{}:
			for key, values := range del {
				for _, value := range values.([]interface{}) {
					for i, v := range mf.Properties[key] {
						if v == value {
							mf.Properties[key] = append(mf.Properties[key][:i], mf.Properties[key][i+1:]...)
						}
					}
				}
			}
		default:
			return fmt.Errorf("unknown delete type %T", del)
		}
	}
	mfData := &microformats.Data{Items: []*microformats.Microformat{mf}}
	entry := mfobjects.GetHEntry(mfData)
	entry.Updated = time.Now().UTC()
	*post = MicropubPost{Entry: entry, RawData: mfData}
	return nil
}
