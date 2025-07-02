package xslt

import(
	"io"
	"os"
	"bytes"
	"errors"
	"os/exec"
)

// XSLTransform applies XSLT rules from file xslFileName to data byte slice
// or data read from inFileName argument if given an uninitialized slice
// and inFileName is not an empty string.
// TODO:Make read from stdin.
// If outFileName is set then output will be written to this file.
// Othervise []byte will be returned.
// The function uses xalan as an XSL transformer.
// TODO: make a possibility to choose a transformer. Put transformer
// arguments some where else.
func XSLTransform(data []byte, inFileName string, xslFileName string, outFileName string) ([]byte, error) {
	//xalan transformation by default
	
	if (data == nil && inFileName == "") || xslFileName == "" {
		return nil, errors.New(ER_XSL_TRANSFORM)
	}	
	
	var out_b []byte
	var errb bytes.Buffer
	params := []string{"-q", "-xsl", xslFileName}
	
	if outFileName != "" {
		params = append(params, "-out", outFileName)
	}	
	if data == nil {
		params = append(params, "-in", inFileName)
	}
	cmd := exec.Command("xalan", params...)
	cmd.Stderr = &errb

	if data != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil { 
			return nil, err
		}

		go func() {
			defer stdin.Close()
			io.Copy(stdin, bytes.NewReader(data))
		}()

		if outFileName != "" {
			err := cmd.Run()
			if err != nil { 
				return nil, err
				//errors.New(string(out_b))
			}		
		}else{
			out_b, err = cmd.Output()		
			if err != nil {
				return nil, err
				//errors.New(string(out_b))
			}
		}

	}else{
		err := cmd.Run()
		if err != nil { 
			return nil, errors.New(string(out_b))
			//errors.New(errb.String())
		}
	}
	
	if outFileName != "" {
		_, err := os.Stat(outFileName)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
			
		}else if err != nil {
			return nil, errors.New(ER_XSL_TRANSFORM+" "+errb.String())
		}
		return nil, nil
	}else{
		return out_b, nil
	}
}

// XSLToPDFTransform transforms given data in form of a byte slice or from inFileName.
// Styles from xslFileName are applied.
// Transformed result is saved to file if outFileName is not empty or is returned
// as slice of bytes.
// The function uses Apache FOP as a transformer. 
// TODO: Make read from stdin.
// This function is used from view/pdf.go
func XSLToPDFTransform(fop string, confFile string, params []string, data []byte, inFileName string, xslFileName string, outFileName string) ([]byte, error) {
	//fop transformation by default	
	if (data == nil && inFileName == "") || xslFileName == "" {
		return nil, errors.New(ER_XSL_TRANSFORM)
	}	
	
	if params == nil || len(params) == 0 {
		params = []string{"-q"} //default param
	}
	if confFile != "" {
		params = append(params, "-c", confFile)
	}
	if fop == "" {
		fop = "fop"
	}
	
	var out_b []byte
	var errb bytes.Buffer
	
	if outFileName != "" {
		params = append(params, "-pdf", outFileName)
	}else{
		params = append(params, "-pdf", "-")
	}	
	if data == nil {
		params = append(params, "-xml", inFileName)
	}else{
		params = append(params, "-xml", "-")
	}
	params = append(params, "-xsl", xslFileName)
	
	//a.Logger.Debugf("XSLToPDFTransform: %s %v\n", fop, params)	
	
	cmd := exec.Command(fop, params...)
	cmd.Stderr = &errb

	if data != nil {
		stdin, err := cmd.StdinPipe()
		if err != nil { 
			return nil, err
		}

		go func() {
			defer stdin.Close()
			io.Copy(stdin, bytes.NewReader(data))
		}()

		if outFileName != "" {
			err := cmd.Run()
			if err != nil { 
				return nil, err
				//errors.New(string(out_b))
			}		
		}else{
			out_b, err = cmd.Output()		
			if err != nil {
				return nil, err
				//errors.New(string(out_b))
			}
		}

	}else{
		err := cmd.Run()
		if err != nil { 
			return nil, errors.New(string(out_b))
		}
	}
	
	if outFileName != "" {
		_, err := os.Stat(outFileName)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
			
		}else if err != nil {
			return nil, errors.New(ER_XSL_TRANSFORM+" "+errb.String())
		}
		return nil, nil
	}else{
		return out_b, nil
	}
}


