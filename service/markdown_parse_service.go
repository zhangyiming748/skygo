package service

import (
	"regexp"
	"strings"
)

const (
	//markdown解析段类型
	MD_SEGMENT_PARA   = "segment_paragraph" //段落
	MD_SEGMENT_HEADER = "segment_header"    //标题

	//markdown解析标签类型
	MD_ELEM_H1           = "h1"
	MD_ELEM_H2           = "h2"
	MD_ELEM_H3           = "h3"
	MD_ELEM_H4           = "h4"
	MD_ELEM_H5           = "h5"
	MD_ELEM_H6           = "h6"
	MD_ELEM_ORDER_LIST   = "order_list"     //有序列表
	MD_ELEM_UNORDER_LIST = "unordered_list" //无序列表

	MD_ELEM_IMG  = "image"
	MD_ELEM_TEXT = "text"
)

var (
	h1Reg          = regexp.MustCompile(`^\s*#\s+([^\n]*)`)
	h2Reg          = regexp.MustCompile(`^\s*##\s+([^\n]*)`)
	h3Reg          = regexp.MustCompile(`^\s*###\s+([^\n]*)`)
	h4Reg          = regexp.MustCompile(`^\s*####\s+([^\n]*)`)
	h5Reg          = regexp.MustCompile(`^\s*#####\s+([^\n]*)`)
	h6Reg          = regexp.MustCompile(`^\s*######\s+([^\n]*)`)
	orderListReg   = regexp.MustCompile(`^(\d{1,10}\.\s[^\n]*)`)
	unorderListReg = regexp.MustCompile(`^-\s([^\n]*)`)
	imgReg         = regexp.MustCompile(`\!\[[^\]]*\]\(([^\)]*)\)`)
)

type MdSegment struct {
	Type     string
	Elements []*MdElem
}

type MdElem struct {
	Type    string
	Content string
}

func ParseMarkdown(mdContent string) []*MdSegment {
	markdown := []*MdSegment{}
	mdSlice := strings.Split(mdContent, "\n")
	for _, mdLine := range mdSlice {
		segmentTag := MD_SEGMENT_HEADER
		elemTag := ""
		phaseSlice := []string{}
		if phaseSlice = h1Reg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_H1
		} else if phaseSlice = h2Reg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_H2
		} else if phaseSlice = h3Reg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_H3
		} else if phaseSlice = h4Reg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_H4
		} else if phaseSlice = h5Reg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_H5
		} else if phaseSlice = h6Reg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_H6
		} else if phaseSlice = orderListReg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_ORDER_LIST
		} else if phaseSlice = unorderListReg.FindStringSubmatch(mdLine); len(phaseSlice) > 0 {
			elemTag = MD_ELEM_UNORDER_LIST
		} else {
			segmentTag = MD_SEGMENT_PARA
		}
		newSegment := &MdSegment{
			Type: segmentTag,
		}
		if segmentTag == MD_SEGMENT_HEADER {
			newElem := &MdElem{
				Type:    elemTag,
				Content: phaseSlice[1],
			}
			newSegment.Elements = []*MdElem{newElem}
		} else {
			newSegment.Elements = parseParagraph(mdLine)
		}
		markdown = append(markdown, newSegment)
	}
	return markdown
}

// 解析段落文字
// 将段落中的文字、图片按照顺序解析出来
func parseParagraph(paraContent string) []*MdElem {
	mdElems := []*MdElem{}
	startIndex := 0
	endIndex := len(paraContent)
	for {
		if startIndex >= endIndex {
			break
		}
		imgIndex := imgReg.FindStringIndex(paraContent[startIndex:endIndex])
		if imgIndex != nil {
			//如果图片标签前面还有文字，则把文字先拿出来
			if imgIndex[0] > 0 {
				textElem := &MdElem{
					Type:    MD_ELEM_TEXT,
					Content: paraContent[startIndex : startIndex+imgIndex[0]],
				}
				mdElems = append(mdElems, textElem)
			}
			if imgMatch := imgReg.FindStringSubmatch(paraContent[startIndex+imgIndex[0] : startIndex+imgIndex[1]]); len(imgMatch) == 2 {
				imgElem := &MdElem{
					Type:    MD_ELEM_IMG,
					Content: imgMatch[1],
				}
				mdElems = append(mdElems, imgElem)
			}
			startIndex += imgIndex[1]
		} else {
			textElem := &MdElem{
				Type:    MD_ELEM_TEXT,
				Content: paraContent[startIndex:endIndex],
			}
			mdElems = append(mdElems, textElem)
			startIndex = endIndex
		}
	}
	return mdElems
}

func IsImage(paraContent string) (string, bool) {
	imgIndex := imgReg.FindStringIndex(paraContent)
	if imgIndex != nil {
		if imgMatch := imgReg.FindStringSubmatch(paraContent[imgIndex[0]:imgIndex[1]]); len(imgMatch) == 2 {
			return imgMatch[1], true
		}
	}
	return "", false
}
